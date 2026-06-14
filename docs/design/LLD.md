# Low-Level Design (LLD) — OSINT Scraper REST API

This document details the concrete interfaces, data structures, and implementation logic for the core subsystems of the OSINT Scraper.

---

## 1. Platform Module & Registry Layer

The platform layer allows developers to add support for new social media providers by implementing a unified interface.

### 1.1 The Platform Interface
Defined in [platform.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/platform/platform.go), the `Platform` interface requires the following structure:

```go
package platform

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Platform interface {
	Name() string
	RegisterRoutes(rg *gin.RouterGroup)
	Ping(ctx context.Context) bool
}
```

Every platform module (e.g. `instagram`, `x`, `snapchat`, `jaco`, `telegram`) implements this interface. 

### 1.2 The Platform Registry
Defined in [registry.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/platform/registry.go), the `Registry` struct acts as a catalog of all active modules:

```go
type Registry struct {
	mu        sync.RWMutex
	platforms []Platform
	log       zerolog.Logger
}
```

* **Dynamic Mounting**: The `MountAll(engine *gin.Engine)` method dynamically groups routes under a prefix based on the platform name: `/api/<platform_name>`.
* **Thread Safety**: Concurrent access to the platforms slice is protected using a `sync.RWMutex`.

---

## 2. Layered Caching Subsystem

The caching layer sits between the incoming API requests and the outbound scraping modules, optimizing performance and reducing external traffic.

### 2.1 The Cache Interface
Defined in [cache.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/cache/cache.go), the `Cache` interface decouples storage logic from routing:

```go
package cache

import "time"

type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, ttl time.Duration)
	Delete(key string)
	Flush()
	Close() error
}
```

### 2.2 MemoryCache (In-Memory implementation)
The `MemoryCache` utilizes a map of `entry` structs protected by a `sync.RWMutex`:

```go
type entry struct {
	data      []byte
	expiresAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]entry
	stop  chan struct{}
}
```

* **Janitor Sweeper**: Spawns a background goroutine via the `janitor` method which runs periodically at the configured `cleanupInterval`. It locks the map and deletes expired keys, preventing memory leaks:
  ```go
  func (c *MemoryCache) evictExpired() {
  	c.mu.Lock()
  	defer c.mu.Unlock()
  	now := time.Now()
  	for k, e := range c.items {
  		if now.After(e.expiresAt) {
  			delete(c.items, k)
  		}
  	}
  }
  ```

### 2.3 RedisCache (Distributed Redis implementation)
Defined in [redis.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/cache/redis.go), the `RedisCache` delegates caching to a Redis instance:

```go
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}
```

* **TTL Management**: Natively relies on Redis expiration. If `err == redis.Nil` occurs, the `Get` method returns `false` to indicate a cache miss.

---

## 3. Core HTTP Client & Service Layers

All outbound requests pass through a standardized HTTP client layer to enforce timeouts, manage cookies, inject rotated headers, and handle errors.

### 3.1 Session Wrapper
Defined in [session.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/httpclient/session.go), the `Session` struct encapsulates an HTTP client and cookies:

```go
type Session struct {
	Client  *http.Client
	Jar     *cookiejar.Jar
	Headers map[string]string
}
```

* **User-Agent & Connection defaults**: Automatically injects realistic `User-Agent` headers (defaulting to standard Chrome/Mac macOS UA) and sets `Connection: keep-alive` to reuse connections.
* **Headers Merging**: Performs a non-destructive merge of default headers, session-level headers (e.g. CSRF tokens, authentication tokens), and per-request headers.

### 3.2 Base Platform Service
Defined in [base_service.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/httpclient/base_service.go), `BaseService` provides the core lifecycle hooks for individual scrapers:

```go
type BaseService struct {
	Session       *Session
	HeaderManager *header.Manager
	Cache         cache.Cache
	Health        *header.HealthTracker
	PlatformName  string
	PlatformURL   string
	Log           zerolog.Logger
}
```

* **Cache Routing**: The `CachedGet` helper intercepts fetch functions to search the cache first and save successful results with dynamic TTL configuration.
* **Rotation Retries (`WithRetry`)**: The core retrying strategy loops up to a maximum attempt limit, querying new rotated headers from the header manager, sending requests, and logging success or failure states to the `HealthTracker`:
  ```go
  func (b *BaseService) WithRetry(ctx context.Context, maxRetries int, fn func(session *Session, headers map[string]string, headerID string) ([]byte, error)) ([]byte, error) {
  	// ...
  	for i := 0; i < maxRetries; i++ {
  		headers, id, err := b.HeaderManager.GetHeaders(b.PlatformName)
  		// ...
  		res, err := fn(session, headers, id)
  		if err == nil {
  			b.Health.ReportSuccess(b.PlatformName, id)
  			return res, nil
  		}
  		b.Health.ReportFailure(b.PlatformName, id)
  	}
  	// ...
  }
  ```

---

## 4. Header Management & Hot Reloading

Rotation and runtime updating of API headers are coordinated through files, watchers, and health state indicators.

### 4.1 Configurations Models
Defined in [model.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/header/model.go), structural configurations mirror the `headers.json` schema:

```go
type HeaderConfig struct {
	Version   string                    `json:"version"`
	Platforms map[string]PlatformConfig `json:"platforms"`
}

type PlatformConfig struct {
	Enabled         bool          `json:"enabled"`
	Rotation        string        `json:"rotation"`
	Headers         []HeaderEntry `json:"headers"`
	CacheTTLSeconds int           `json:"cache_ttl_seconds,omitempty"`
}

type HeaderEntry struct {
	ID      string            `json:"id"`
	Headers map[string]string `json:"headers"`
	Enabled *bool             `json:"enabled,omitempty"`
}
```

### 4.2 Header Manager
Defined in [manager.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/header/manager.go), the `Manager` regulates platforms, active configurations, and rotation states:

```go
type Manager struct {
	cfg      *HeaderConfig
	counters map[string]*counter
	mu       sync.RWMutex
}
```

* **Round-Robin Rotation**: The `Next` function resolves round-robin indexing against the number of *currently enabled* header entries for a specific platform. If no headers are enabled, it returns an error.
* **Locking**: Employs write-locks (`m.mu.Lock()`) for CRUD operations, and read-locks (`m.mu.RLock()`) for config queries and TTL resolutions to ensure complete data consistency across concurrent requests.

### 4.3 fsnotify Hot Reloading
Defined in [watcher.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/header/watcher.go), `WatchAndReload` listens to write and creation events on the JSON file:

* When a filesystem change is received, the file is re-parsed via `LoadHeaderConfig`.
* On success, `Manager.Reload` is executed, hot-swapping the active configuration in memory without drop-outs or restarts.

### 4.4 HealthTracker Auto-Disabling
Defined in [health.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/header/health.go), the circuit breaker monitors failure metrics per header entry:

```go
type HealthTracker struct {
	mu        sync.Mutex
	stats     map[string]*EntryStats
	threshold int
	manager   *Manager
	log       zerolog.Logger
}
```

* **Fails Threshold**: When a header entry's consecutive error count exceeds `threshold` (default is **5**), `HealthTracker` automatically calls `Manager.UpdateHeader` to disable that key.
* **Auto-Recovery Reset**: When a request succeeds, the `ConsecutiveFails` counter for the matched header entry is reset to `0`.
