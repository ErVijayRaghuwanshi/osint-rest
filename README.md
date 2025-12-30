# OSINT Scraper REST API

A high-performance, scalable, and extensible **OSINT scraping microservice** built with **Go** and **Gin**, featuring **Swagger (OpenAPI) documentation** and a **clean, modular architecture**.

This service is designed to support **multiple social media platforms** (Instagram, X/Twitter, Snapchat, etc.) with isolated modules, making it easy to add new platforms, scraping strategies, or data pipelines.

---

## Features

- ⚡ **Gin-based REST API** for high throughput
- 📘 **Swagger UI** for interactive API documentation
- 🧩 **Modular platform-based architecture**
  - Instagram
  - X (Twitter)
  - Snapchat
- 🧠 Clear separation of concerns:
  - Routes
  - Handlers
  - Services
  - Models
- 🐳 **Docker-ready**, Kubernetes-friendly
- 🛡 Designed for OSINT use-cases:
  - Proxy rotation (future)
  - Session & header management
  - Rate limiting & retries
  - Concurrent scraping workers

---

## Project Structure

```

osint-scraper/
├── cmd/
│   └── api/
│       └── main.go                # Application entrypoint
│
├── internal/
│   ├── api/
│   │   ├── handlers.go            # Generic endpoints (health, base)
│   │   ├── router.go              # Central API router
│   │   │
│   │   ├── instagram/             # Instagram module
│   │   │   ├── handlers.go
│   │   │   ├── routes.go
│   │   │   ├── service.go
│   │   │   └── models.go
│   │   │
│   │   ├── snapchat/              # Snapchat module
│   │   │   ├── handlers.go
│   │   │   ├── routes.go
│   │   │   └── service.go
│   │   │
│   │   └── x/                     # X (Twitter) module
│   │       ├── handlers.go
│   │       ├── routes.go
│   │       ├── service.go
│   │       └── models.go
│   │
│   ├── config/
│   │   └── config.go               # Application configuration
│   │
│   └── logger/
│       └── logger.go               # Structured logging
│
├── docs/                           # Swagger / OpenAPI docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
├── README.md
└── Testing.ipynb                   # Research & exploration (non-prod)

````

---

## Prerequisites

- **Go 1.22+**
- **Swag CLI** (for Swagger docs)
- Optional:
  - Docker
  - Kubernetes

---

## Installation

### 1️⃣ Clone Repository

```bash
git clone https://github.com/ErVijayRaghuwanshi/osint-scraper.git
cd osint-scraper
````

### 2️⃣ Install Dependencies

```bash
make deps
make tidy
```

### 3️⃣ Generate Swagger Docs

```bash
make swagger
```

---

## Running the API

```bash
make run
```

* API runs on **port 8080** (configurable via env)
* Swagger UI:

  ```
  http://localhost:8080/swagger/index.html
  ```

---

## Example Endpoints

### Health Check

```http
GET /health
```

```json
{
  "status": "ok"
}
```

---

### Snapchat Ping

```http
GET /api/snapchat/ping
```

```json
{
  "message": "snapchat pong"
}
```

---

### Instagram Ping

```http
GET /api/instagram/ping
```

---

### X (Twitter) Ping

```http
GET /api/x/ping
```

---

## Adding a New Platform Module

Example: **Telegram**

### 1️⃣ Create module directory

```
internal/api/telegram/
├── handlers.go
├── routes.go
├── service.go
├── models.go
```

### 2️⃣ Register routes

```go
telegram.RegisterRoutes(router.Group("/api/telegram"))
```

### 3️⃣ Add Swagger annotations in handlers

---

## Docker Usage

### Build Image

```bash
make docker-build
```

### Run Container

```bash
make docker-run
```

---

## Kubernetes

The service is **Kubernetes-ready**.

Typical resources:

* Deployment
* Service
* ConfigMap / Secrets
* Ingress (optional)

```bash
make k8s-apply
```

---

## Future Enhancements

* Proxy rotation per platform
* Cookie & session pools
* Distributed scraping workers
* Elasticsearch integration
* Async task queue (Redis / Kafka)
* Auth (JWT / API keys)
* Rate limiting & abuse protection

---

## License

MIT License © 2025

---

## Author

**Er Vijay Raghuwanshi**
📧 Email: [ervijayraghuwanshi@gmail.com](mailto:ervijayraghuwanshi@gmail.com)

---

> Built with ❤️ for OSINT, threat intelligence, and scalable backend systems.

````
