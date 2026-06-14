# OSINT Scraper — Feature Guide & Operations

This document details the inner workings of the Caching System, the Admin Control APIs, and the Header Management subsystem (including reloading, cooldown, and proxy configurations).

---

## 1. Caching System Architecture

The OSINT Scraper features a **layered, non-intrusive caching architecture** that intercepts requests at the HTTP routing level. This design decouples caching concerns from the platform scraping packages.

```
Incoming Request (GET)
         │
         ▼
┌──────────────────┐
│ CacheMiddleware  │
└────────┬─────────┘
         │
         ├──► [Exists in Cache?] ──(Yes)──► Return 200 OK (X-Cache: HIT)
         │
         └──► (No) ──► Forward to Scraper Service ──► Outbound API request
                            │
                            ▼
                     [Response 200 OK?] ──(Yes)──► Save to Cache (X-Cache: MISS)
```

### 1.1 Layered HTTP Middleware
The `CacheMiddleware` (implemented in [cache.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/api/middleware/cache.go)) wraps Gin's `ResponseWriter` using a custom buffer writer. 

* **Cache Keys**: Keys are formatted as `http:cache:<RequestURI>` (e.g. `http:cache:/api/jaco/userinfo?username=Uaegirl`). Since the full Request URI is used, query parameters (such as `username` or `limit`) are naturally isolated in different cache entries.
* **Dynamic TTL Resolution**: The middleware parses the platform identifier from the request path (e.g., `/api/instagram/userinfo` -> `instagram`) and calls `hm.GetCacheTTL(platform)`.
  * If a TTL is defined for the platform in `headers.json` (via `cache_ttl_seconds`), it is applied.
  * If not configured, it falls back to a default TTL of **5 minutes**.

### 1.2 Configurable Storage Drivers
The caching layer implements a generic `cache.Cache` interface, allowing you to configure the backend via the `CACHE_TYPE` environment variable in your `.env`:

#### Option A: In-Memory Cache (`CACHE_TYPE=memory`)
* Uses Go's native map protected by a `sync.RWMutex`.
* Boots a background "janitor" goroutine that sweeps and evicts expired cache entries at a configured interval (default: every 2 minutes).
* Recommended for single-instance or local development.

#### Option B: Redis Cache (`CACHE_TYPE=redis`)
* Instantiates a Redis client using a connection pool configured by `REDIS_ADDR`, `REDIS_PASSWORD`, and `REDIS_DB`.
* Natively delegates key expiration and memory eviction to the Redis engine.
* Recommended for distributed deployments to support horizontal scaling across multiple container replicas.

---

## 2. Admin Control API

The Admin API provides remote management of rotated HTTP headers, authorization tokens, and active scraping platforms at runtime.

### 2.1 Security & Authentication
All endpoints under the `/admin` path group are protected by the `AdminAPIKey` middleware:
* It checks for the presence of the `X-API-Key` HTTP request header.
* The header value must match the `ADMIN_API_KEY` environment variable configured in `.env`.
* If the key is missing or incorrect, it immediately terminates the request returning `401 Unauthorized`.

### 2.2 Endpoint Registry & Payloads

| HTTP Method | Route | Description | Request Payload |
|---|---|---|---|
| **`GET`** | `/admin/headers` | List all platforms and their rotation configs | *None* |
| **`GET`** | `/admin/headers/:platform` | List all header entries for a specific platform | *None* |
| **`POST`** | `/admin/headers/:platform` | Add a new header entry | JSON Header Entry |
| **`PUT`** | `/admin/headers/:platform/:id` | Update an existing header entry | JSON Header Entry |
| **`DELETE`** | `/admin/headers/:platform/:id` | Delete a header entry | *None* |
| **`POST`** | `/admin/headers/:platform/:id/toggle` | Enable or disable a specific header entry | *None* |

### 2.3 Persistence Model
When a modifying request is sent (e.g. `POST`, `PUT`, `DELETE`), the Admin Handlers (in [handlers.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/api/admin/handlers.go)) update the `HeaderManager` memory map and immediately write the updated JSON configuration back to the local file `/config/headers.json`.

---

## 3. Advanced Outbound request Dynamics

To ensure robust data collection and bypass anti-bot challenges, the scraper incorporates hot-reloading, failure-aware cooldowns, and proxy capabilities.

### 3.1 Header Reloading (Zero-Downtime Hot Reload)
The `Watcher` engine (implemented in [watcher.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/header/watcher.go)) utilizes `fsnotify` to listen for filesystem modification write events on `config/headers.json`.
* When the file is updated (either by the Admin API or via a Kubernetes ConfigMap/Persistent Volume update), the file watcher detects the write.
* It immediately calls `HeaderManager.Reload()`, parsing the new JSON structure and updating the active headers memory maps.
* **Benefit**: Credentials, sessions, and User-Agents can be updated dynamically without restarting the server or dropping active connections.

### 3.2 Header Cooldown / Auto-Disabling
The `HealthTracker` (implemented in [health.go](file:///Users/ervijay/Documents/Programs/Repo/osint-scraper/internal/header/health.go)) acts as a circuit breaker for individual header sets:
* **Failure Count**: If an outbound request to a target platform fails (e.g., returns `403 Forbidden` or `429 Too Many Requests`), the platform service reports a failure.
* **Auto-Disable Trigger**: If a specific header entry (e.g. `insta-1`) accumulates **5 consecutive failures**, the Health Tracker disables the header (`HeaderManager.ToggleHeaderSet(platform, id, false)`) and writes this state to the JSON config.
* **Benefit**: Keeps the rotation clean by removing rate-limited or expired cookies from the active rotation pool, allowing them to "cool down" until updated by an administrator.

### 3.3 Proxy Integration Strategy
While standard requests use direct outbound routing, the shared HTTP client is designed for proxy rotation:
* **Roadmap Design**: Future phases map a `proxy` configuration field to each header entry in `headers.json`:
  ```json
  {
    "id": "insta-1",
    "headers": { ... },
    "proxy": "http://user:pass@residential-proxy-ip:port"
  }
  ```
* **Pairing**: In `httpclient.Session`, when initializing an outbound request, a custom `http.Transport` will dial through the paired proxy. This pairs specific cookies/sessions with fixed residential IPs, making scraper requests appear highly realistic to third-party anti-bot checkers.
