# OSINT Scraper REST API

A high-performance, scalable, and extensible **OSINT scraping microservice** built with **Go** and **Gin**, featuring **Swagger (OpenAPI) documentation** and a **clean, modular architecture**.

This service is designed to support **multiple social media platforms** (Instagram, X/Twitter, Snapchat, Jaco) with isolated modules, making it easy to add new platforms, scraping strategies, or data pipelines.

---

## Features

- **Gin-based REST API** for high throughput
- **Swagger UI** for interactive API documentation
- **Modular platform-based architecture**
  - Instagram
  - X (Twitter)
  - Snapchat
  - Jaco
- **Shared HTTP client** with session management, cookie jars, and header injection
- **In-memory response cache** with per-platform TTL configuration
- **Header management system**
  - Round-robin rotation of auth headers per platform
  - Hot reload via `fsnotify` (no restart needed)
  - Admin REST API for CRUD operations on headers (API key protected)
  - Health scoring with auto-disable of failing headers
- **Production hardened**
  - `context.Context` propagation for cancellable requests
  - Graceful shutdown (SIGINT/SIGTERM)
  - Per-IP rate limiting
  - CORS middleware
  - Structured logging via `zerolog`
- **Platform interface & registry** for expandable platform integration
- **Docker-ready**, Kubernetes-friendly
- **Unit tests** for header manager, cache, and health tracker

---

## Project Structure

```
osint-scraper/
├── cmd/api/
│   └── main.go                     # Application entrypoint
│
├── internal/
│   ├── api/
│   │   ├── handlers.go             # Health check endpoint
│   │   ├── router.go               # Central API router
│   │   ├── admin/                  # Admin API (header CRUD)
│   │   │   ├── handlers.go
│   │   │   └── routes.go
│   │   ├── middleware/             # CORS, rate limiting, API key auth
│   │   │   ├── apikey.go
│   │   │   ├── cors.go
│   │   │   └── ratelimit.go
│   │   ├── common/                 # Shared response types
│   │   │   └── responses.go
│   │   ├── instagram/              # Instagram module
│   │   ├── snapchat/               # Snapchat module
│   │   ├── x/                      # X (Twitter) module
│   │   └── jaco/                   # Jaco module
│   │
│   ├── httpclient/                 # Shared HTTP session & base service
│   │   ├── session.go
│   │   └── base_service.go
│   │
│   ├── header/                     # Header management system
│   │   ├── manager.go              # Core logic + CRUD
│   │   ├── model.go                # JSON structs
│   │   ├── loader.go               # File loading
│   │   ├── roundrobin.go           # Rotation logic
│   │   ├── watcher.go              # fsnotify hot reload
│   │   └── health.go               # Health scoring & auto-disable
│   │
│   ├── cache/                      # In-memory cache with TTL
│   │   └── cache.go
│   │
│   ├── platform/                   # Platform interface & registry
│   │   ├── platform.go
│   │   └── registry.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   └── logger/
│       └── logger.go
│
├── config/
│   ├── headers.json                # Header config (gitignored)
│   └── headers.example.json        # Template
│
├── k8s/                            # Kubernetes manifests
├── docs/                           # Swagger / OpenAPI docs
├── Dockerfile
├── Makefile
├── .env.example
├── .gitignore
└── README.md
```

---

## Prerequisites

- **Go 1.22+**
- **Swag CLI** (for Swagger docs)
- Optional:
  - Docker
  - Kubernetes

---

## Installation

### 1. Clone Repository

```bash
git clone https://github.com/ErVijayRaghuwanshi/osint-scraper.git
cd osint-scraper
```

### 2. Setup Configuration

```bash
cp .env.example .env
cp config/headers.example.json config/headers.json
# Edit both files with your actual values
```

### 3. Install Dependencies

```bash
make deps
make tidy
```

### 4. Generate Swagger Docs

```bash
make swagger
```

---

## Running the API

```bash
make run
```

- API runs on **port 8080** (configurable via `SERVER_PORT` env)
- Swagger UI: `http://localhost:8080/swagger/index.html`

---

## Example Endpoints

### Health Check

```http
GET /health
```

```json
{
  "status": "ok",
  "goroutines": 8,
  "alloc_mb": 2,
  "go_version": "go1.25.5"
}
```

### Platform Ping

```http
GET /api/snapchat/ping
GET /api/instagram/ping
GET /api/x/ping
GET /api/jaco/ping
```

### User Info

```http
GET /api/snapchat/userinfo?username=arora_girl
GET /api/instagram/userinfo?username=sakshi_raghu_1c_
GET /api/x/userinfo?screen_name=urstrulymahesh
GET /api/jaco/userinfo?username=Uaegirl
```

### Admin API (requires `X-API-Key` header)

```http
GET    /admin/headers                     # List all platforms
GET    /admin/headers/:platform           # List headers for platform
POST   /admin/headers/:platform           # Add header set
PUT    /admin/headers/:platform/:id       # Update header set
DELETE /admin/headers/:platform/:id       # Delete header set
POST   /admin/headers/:platform/:id/toggle # Toggle enabled/disabled
```

---

## Adding a New Platform Module

1. Create module directory:
   ```
   internal/api/telegram/
   ├── handlers.go   (embed zerolog.Logger, use c.Request.Context())
   ├── routes.go
   ├── service.go    (embed httpclient.BaseService)
   └── models.go
   ```

2. Register routes in `internal/api/router.go`

3. Add platform headers to `config/headers.json`

4. Add Swagger annotations in handlers

---

## Docker Usage

```bash
make docker-build
make docker-run
```

---

## Kubernetes

```bash
make k8s-apply
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | API server port |
| `SWAGGER_HOST` | `localhost:8080` | Swagger UI host |
| `GIN_MODE` | `release` | Gin mode (debug/release) |
| `HEADERS_CONFIG_PATH` | `./config/headers.json` | Path to header config |
| `ADMIN_API_KEY` | (none) | API key for admin endpoints |
| `API_TITLE` | `OSINT Scraper API` | Swagger title |
| `API_DESCRIPTION` | (default) | Swagger description |

---

## Future Enhancements

- Proxy rotation per platform
- Redis cache adapter for distributed deployments
- Weighted & failover header rotation strategies
- Platform scaffold generator (`make new-platform name=telegram`)
- Elasticsearch integration
- Async task queue (Redis / Kafka)
- Metrics & observability (Prometheus)

---

## License

MIT License

---

## Author

**Er Vijay Raghuwanshi**
Email: [ervijayraghuwanshi@gmail.com](mailto:ervijayraghuwanshi@gmail.com)

---

> Built for OSINT, threat intelligence, and scalable backend systems.
