# Future Architecture Roadmap — OSINT Scraper

This document outlines the roadmap for the OSINT Scraper service, broken down into sequential execution phases.

---

## Phase 1: Advanced Proxy & User-Agent (UA) Rotation

**Objective**: Maximize evasion of bot detectors and rate-limit blocks by routing outbound scraper requests through proxy networks.

### 1. Proxy-to-Header Pairing
* Add a `proxy` field to individual `HeaderEntry` items in [model.go](../../internal/header/model.go):
  ```json
  {
    "id": "instagram-session-1",
    "headers": {
      "Cookie": "sessionid=...",
      "User-Agent": "..."
    },
    "proxy": "socks5://user:pass@12.34.56.78:1080"
  }
  ```
* Ensure specific sessions are pinned to the same residential proxy IP. This prevents anti-bot systems from flagging account cookies when requests arrive from different geographic regions in short succession.

### 2. Dynamically Rotated HTTP Transports
* Extend [session.go](../../internal/httpclient/session.go) to instantiate a dedicated `http.Transport` for each distinct proxy configuration.
* Enable HTTP/2 and HTTP/3 support on proxy connections to mirror browser behaviors.

### 3. User-Agent Spoofing Engine
* Integrate a dynamic UA generation library to automatically generate matching header signatures (e.g. matching `Sec-Ch-Ua` and platform capabilities with the client's User-Agent string).

---

## Phase 2: Distributed State & Messaging

**Objective**: Scale the service horizontally to support high-throughput, multi-replica deployments.

```
                  ┌──────────────────────┐
                  │ OSINT Scraper Node 1 │
                  └──────────┬───────────┘
                             │
                             ▼ (Distributed Locks / Cache)
    ┌──────────────────────────────────────────────────┐
    │                  Redis Cluster                   │
    └──────────────────────────────────────────────────┘
                             ▲ (Distributed Locks / Cache)
                             │
                  ┌──────────┴───────────┐
                  │ OSINT Scraper Node 2 │
                  └──────────────────────┘
```

### 1. Redis Distributed Caching & Rate Limiting
* Fully transition default deployments from `MemoryCache` to [redis.go](../../internal/cache/redis.go).
* Implement distributed rate-limiting middleware using Redis sliding-window counters to ensure global request thresholds are respected across multiple instances.

### 2. Distributed Lock Management (Redlock)
* Use Redis-based distributed locking (e.g. `go-redlock`) during platform updates and configuration reloads.
* Prevent "cache stampedes" where multiple scraper nodes simultaneously request the same user data on cache miss.

### 3. Asynchronous Task Queue (Asynq / RabbitMQ)
* Integrate `Asynq` or RabbitMQ for long-running scraping tasks (e.g. fetching entire user post histories).
* Provide webhook callbacks to notify upstream clients once data processing completes.

---

## Phase 3: Telemetry, Credentials Vaults & Client Auth

**Objective**: Improve operational visibility, secure sensitive API credentials, and regulate API access.

### 1. Prometheus Metrics Exporter
* Export real-time metrics including:
  * Cache hit/miss rates.
  * Platform scraping latencies (quantiles).
  * Outbound HTTP status codes (2xx, 4xx, 5xx) per target platform.
  * Active/disabled header entries count.
* Package default Grafana dashboards for quick visualization of scraping health.

### 2. HashiCorp Vault Integration
* Move sensitive session cookies and API keys out of [headers.json](../../config/headers.json) and retrieve them dynamically from HashiCorp Vault.
* Support hot reloading of credentials directly from Vault transit paths.

### 3. Client Authentication & Tenant Isolation
* Introduce OAuth2/OIDC validation for incoming requests to the scraper.
* Support multi-tenant rate limits and usage quotas.

---

## Phase 4: React / Next.js Admin Dashboard

**Objective**: Build a premium, user-friendly frontend dashboard to replace the command-line/REST API control of headers and rotation pools.

```
┌────────────────────────────────────────────────────────────────────────────┐
│ OSINT Scraper Admin Console                                                 │
├────────────────────────────────────────────────────────────────────────────┤
│ Platforms:                                                                 │
│   Instagram [|||||||||||||||||| 100%] Active (3/3 headers healthy)         │
│   Snapchat  [||||||||||.......  60%] Warning (1/2 headers disabled)        │
├────────────────────────────────────────────────────────────────────────────┤
│ Active Header Pools:                                                       │
│   [Add Header Entry] [Reload Config] [Flush Cache]                         │
│                                                                            │
│   ID          Platform    Status      Successes   Failures    Actions      │
│   insta-1     instagram   [HEALTHY]   1,204       0           [Disable]    │
│   insta-2     instagram   [HEALTHY]   945         2           [Disable]    │
│   snap-1      snapchat    [COOLDOWN]  312         5           [Enable]     │
└────────────────────────────────────────────────────────────────────────────┘
```

### 1. Real-time Status Monitoring
* Build a React/Next.js dashboard utilizing WebSockets to show real-time stats from the [HealthTracker](../../internal/header/health.go).
* Display live success/failure rates and highlight disabled headers.

### 2. Pool Management UI
* Provide form controls to add, update, delete, and toggle headers dynamically.
* Implement file upload features to bulk-import header JSON configurations.
