# High-Level Design (HLD) — OSINT Scraper REST API

This document describes the high-level architecture, subsystem boundaries, data flows, and infrastructure models for the OSINT Scraper service.

---

## 1. System Context & Boundaries

The OSINT Scraper acts as a specialized proxy and API layer between internal analysis tools (threat intelligence platforms, investigation UIs) and target social media platforms. It abstracts away session management, header rotation, rate-limiting bypasses, and scraping mechanics.

```mermaid
graph LR
    AnalystUI["Analyst UI / Cron Tasks"] -->|"HTTP REST API"| Scraper["OSINT Scraper Service"]
    Scraper -->|"Outbound HTTP - Rotated"| Snapchat["Snapchat API"]
    Scraper -->|"Outbound HTTP - Rotated"| Instagram["Instagram API"]
    Scraper -->|"Outbound HTTP - Rotated"| Twitter["X / Twitter API"]
    Scraper -->|"Outbound HTTP - Rotated"| Jaco["Jaco API"]
```

---

## 2. High-Level Subsystems

The service is divided into three primary layers:

```
┌────────────────────────────────────────────────────────┐
│                      API Layer                         │
│ - Gin HTTP Router & Middleware (Rate limit, CORS, Auth)│
│ - Swagger UI Specs                                     │
└──────────────────────────┬─────────────────────────────┘
                           │
                           ▼
┌────────────────────────────────────────────────────────┐
│                   Registry Layer                       │
│ - Dynamic route registration & Discovery               │
└──────────────────────────┬─────────────────────────────┘
                           │
                           ▼
┌────────────────────────────────────────────────────────┐
│                   Platform Engines                     │
│ - Snapchat, Instagram, X, Jaco, Telegram Services      │
└──────────────────────────┬─────────────────────────────┘
                           │
                           ▼
┌────────────────────────────────────────────────────────┐
│                   Core Client Engine                   │
│ - Session Wrapper, Cache, Header Rotation & Watcher    │
└────────────────────────────────────────────────────────┘
```

1. **API Layer**: Handles incoming HTTP requests, CORS, rate limiting, Swagger spec generation, and overall lifecycle.
2. **Platform Registry Layer**: Discovers available scrapers, mounts their routes dynamically, and exposes health state.
3. **Platform Engines**: Implement platform-specific parsing, query parameters, API endpoints, and payload formats.
4. **Core Client Engine**: Provides shared services for HTTP sessions, cookie storage, in-memory caching, header file watching, and automated rotation.

---

## 3. Core Architecture Flows

### 3.1 Scraper Request Flow
The following sequence diagram demonstrates the flow of a standard user request (e.g., retrieving Instagram user info):

```mermaid
sequenceDiagram
    autonumber
    actor User as Client / User
    participant Router as Gin Router
    participant Limiter as Rate Limiter
    participant Handler as Instagram Handler
    participant Svc as Instagram Service
    participant Cache as Cache
    participant HM as Header Manager
    participant Client as HTTP Client
    participant Provider as Instagram API
    participant Tracker as Health Tracker

    User->>Router: GET /api/instagram/userinfo?username=foo
    Router->>Limiter: Check Limit
    Limiter-->>Router: Limit OK
    Router->>Handler: Dispatch Request
    Handler->>Svc: GetUserInfo("foo")
    Svc->>Cache: Get("instagram:userinfo:foo")
    alt Cache Hit
        Cache-->>Svc: Return Cached JSON
    else Cache Miss
        Svc->>HM: GetHeaders("instagram")
        HM-->>Svc: Return (Headers, HeaderID)
        Svc->>Client: Outbound GET with Headers
        Client->>Provider: Outbound API Call
        Provider-->>Client: Response (200 OK)
        Client-->>Svc: Raw Body
        Svc->>Cache: Set("instagram:userinfo:foo", Body)
        Svc->>Tracker: ReportSuccess("instagram", HeaderID)
    end
    Svc-->>Handler: Parse response data
    Handler-->>Router: Return 200 OK with Data
    Router-->>User: JSON Response
```

### 3.2 Header Rotation & Error Handling Flow
When a platform request encounters failures (e.g., rate-limiting `429` or forbidden `403` responses), the system automatically retries with a new header set and reports failures:

```mermaid
sequenceDiagram
    autonumber
    participant Svc as Platform Service
    participant Client as HTTP Client
    participant Provider as External API
    participant Tracker as Health Tracker
    participant HM as Header Manager

    Svc->>Client: Request Attempt 1 (Header A)
    Client->>Provider: Send HTTP Request
    Provider-->>Client: 403 Forbidden
    Client->>Tracker: ReportFailure(Header A)
    Note over Tracker: Count fails for Header A
    alt Failures >= 5
        Tracker->>HM: Disable Header A
    end
    Svc->>HM: GetHeaders (Next)
    HM-->>Svc: Return Header B
    Svc->>Client: Request Attempt 2 (Header B)
    Client->>Provider: Send HTTP Request
    Provider-->>Client: 200 OK
    Client->>Tracker: ReportSuccess(Header B)
    Note over Tracker: Reset fail counter for Header B
```

### 3.3 Dynamic Configuration Reloading
The Header Manager watches `config/headers.json` on disk. When edited (e.g., via the Admin API or direct Kubernetes mount updates), the watcher reloads it into memory without server restarts.

```mermaid
graph TD
    Watcher[File Watcher fsnotify] -->|Detects Write| Reload[Reload Configuration]
    Reload -->|Parse JSON| Manager[Header Manager Memory]
    Manager -->|Update maps & counters| Core[Active Scrapers]
```

---

## 4. Infrastructure & Deployment Model

The microservice is designed for containerized deployment in Kubernetes.

* **ConfigMap / PersistentVolume Mounts**: The active `headers.json` file is mounted into the container at `/config/headers.json`. Using a `PersistentVolumeClaim` allows read/write access so the Admin API can persist runtime header modifications.
* **Graceful Shutdown**: Upon receiving `SIGINT` or `SIGTERM`, the Gin server stops accepting new connections and waits for active requests to finish (with a 10-second timeout) before cleaning up the cache and shutting down.
