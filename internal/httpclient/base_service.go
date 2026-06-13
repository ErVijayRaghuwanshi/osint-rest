package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"osint-scraper/internal/cache"
	"osint-scraper/internal/header"

	"github.com/rs/zerolog"
)

// BaseService provides shared session management, header rotation,
// caching, health tracking, and website-reachability logic that every
// platform service needs.
type BaseService struct {
	Session       *Session
	HeaderManager *header.Manager
	Cache         cache.Cache
	Health        *header.HealthTracker
	PlatformName  string
	PlatformURL   string // e.g. "https://www.instagram.com"
	Log           zerolog.Logger
}

// SetHealthTracker assigns a shared health tracker instance.
func (b *BaseService) SetHealthTracker(ht *header.HealthTracker) {
	b.Health = ht
}

// NewBaseService creates a BaseService for the given platform.
func NewBaseService(platformName, platformURL string, hm *header.Manager, log zerolog.Logger) *BaseService {
	return &BaseService{
		HeaderManager: hm,
		PlatformName:  platformName,
		PlatformURL:   platformURL,
		Log:           log,
	}
}

// SetCache assigns a shared cache instance to the service.
func (b *BaseService) SetCache(c cache.Cache) {
	b.Cache = c
}

// CacheTTL returns the configured cache TTL for this platform.
func (b *BaseService) CacheTTL() time.Duration {
	if b.HeaderManager == nil {
		return 0
	}
	secs := b.HeaderManager.GetCacheTTL(b.PlatformName)
	if secs <= 0 {
		return 5 * time.Minute // default
	}
	return time.Duration(secs) * time.Second
}

// CacheResult wraps a cached response with metadata about whether it was
// served from cache, so handlers can set response headers like X-Cache.
type CacheResult struct {
	Data     []byte
	CacheHit bool
}

// CachedGet checks the cache for key, and if missing, calls fetch(),
// stores the result, and returns a CacheResult indicating hit/miss.
func (b *BaseService) CachedGet(key string, fetch func() ([]byte, error)) (*CacheResult, error) {
	if b.Cache != nil {
		if data, ok := b.Cache.Get(key); ok {
			b.Log.Debug().Str("key", key).Msg("Cache hit")
			return &CacheResult{Data: data, CacheHit: true}, nil
		}
	}

	data, err := fetch()
	if err != nil {
		return nil, err
	}

	if b.Cache != nil {
		b.Cache.Set(key, data, b.CacheTTL())
		b.Log.Debug().Str("key", key).Dur("ttl", b.CacheTTL()).Msg("Cache set")
	}

	return &CacheResult{Data: data, CacheHit: false}, nil
}

// InitializeSession creates a new Session and loads headers from the manager.
func (b *BaseService) InitializeSession() (*Session, error) {
	session, err := NewSession()
	if err != nil {
		return nil, err
	}

	if b.HeaderManager != nil {
		headers, id, err := b.HeaderManager.GetHeaders(b.PlatformName)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s headers: %w", b.PlatformName, err)
		}
		b.Log.Info().Str("header_id", id).Msg("Loaded headers for session")
		session.SetAuthHeaders(headers)
	}

	b.Session = session
	return session, nil
}

// EnsureSession returns the current session or initializes a new one.
func (b *BaseService) EnsureSession() (*Session, error) {
	if b.Session != nil {
		return b.Session, nil
	}
	return b.InitializeSession()
}

// CheckWebsite checks if the platform's base URL is reachable.
func (b *BaseService) CheckWebsite(ctx context.Context) bool {
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", b.PlatformURL, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// WithRetry executes fn up to maxRetries times, rotating headers on each attempt.
// It is useful for scraping endpoints that may return 403 with stale headers.
func (b *BaseService) WithRetry(ctx context.Context, maxRetries int, fn func(session *Session, headers map[string]string, headerID string) ([]byte, error)) ([]byte, error) {
	if b.HeaderManager == nil {
		return nil, fmt.Errorf("header manager not set for %s", b.PlatformName)
	}

	var lastErr error

	for i := 0; i < maxRetries; i++ {
		headers, id, err := b.HeaderManager.GetHeaders(b.PlatformName)
		if err != nil {
			lastErr = err
			continue
		}

		b.Log.Debug().
			Int("attempt", i+1).
			Str("header_id", id).
			Msg("Retry attempt")

		session, err := b.EnsureSession()
		if err != nil {
			lastErr = err
			continue
		}

		res, err := fn(session, headers, id)
		if err == nil {
			if b.Health != nil {
				b.Health.ReportSuccess(b.PlatformName, id)
			}
			return res, nil
		}

		if b.Health != nil {
			b.Health.ReportFailure(b.PlatformName, id)
		}
		lastErr = err
	}

	return nil, fmt.Errorf("%s failed after %d retries: %w", b.PlatformName, maxRetries, lastErr)
}
