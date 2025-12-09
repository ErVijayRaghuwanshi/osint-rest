# OSINT Scraper REST API

A high-performance, scalable, and extensible OSINT scraping microservice built with **Go** and **Gin**, featuring **Swagger UI documentation**.  
This project is designed to easily add social media modules (Snapchat, Twitter, Instagram, etc.) with endpoints that check site availability or perform scraping tasks.

---

## Features

- **Gin-based REST API** for high throughput
- **Swagger UI** for API documentation
- **Modular design** for social media services
- **Snapchat endpoint**: `/api/snapchat/ping`
- **Extensible architecture** for:
  - Proxy rotation
  - Session management (cookies, headers)
  - Worker pools for concurrent scraping
- Ready for **Docker & Kubernetes deployment**

---

## Project Structure

```

osint-scraper/
├── cmd/
│   └── api/
│       └── main.go           # Entry point
├── internal/
│   ├── api/
│   │   ├── handlers.go       # Generic endpoints (health, ping)
│   │   ├── router.go         # API routing
│   │   └── snapchat/         # Snapchat module
│   │       ├── handler.go
│   │       ├── routes.go
│   │       └── service.go
│   ├── config/
│   └── logger/
├── docs/                     # Swagger documentation (auto-generated)
├── go.mod
└── go.sum

````

---

## Prerequisites

- Go 1.22+ installed
- `swag` CLI for Swagger documentation
- Optional: Docker for containerized deployment

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/ErVijayRaghuwanshi/osint-scraper.git
cd osint-scraper
````

2. Initialize Go modules:

```bash
go mod tidy
```

3. Install Swag CLI (if not installed):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

4. Generate Swagger documentation:

```bash
swag init -g cmd/api/main.go -o docs
```

---

## Running the API

```bash
go run cmd/api/main.go
```

* API server will start on default port `8080` (configurable via `PORT` environment variable)
* Swagger UI available at: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## Example Endpoints

### Health Check

```http
GET /health
```

**Response**

```json
{
  "status": "ok"
}
```

### Snapchat Ping

```http
GET /api/snapchat/ping
```

**Response (reachable)**

```json
{
  "message": "snapchat pong"
}
```

**Response (unreachable)**

```json
{
  "message": "snapchat unreachable"
}
```

---

## Adding a New Social Media Module

1. Create a folder under `internal/api/`, e.g., `twitter/`
2. Add:

   * `service.go` → business logic
   * `handler.go` → Gin HTTP handler
   * `routes.go` → register routes to router
3. Register the routes in `router.go`:

```go
twitterService := twitter.NewService()
twitter.RegisterRoutes(r.Group("/api/twitter"), twitterService)
```

4. Add Swagger annotations in the handler for documentation.

---

## Docker Deployment

### Build Docker Image

```bash
docker build -t osint-scraper:latest .
```

### Run Docker Container

```bash
docker run -p 8080:8080 osint-scraper:latest
```

---

## Kubernetes Deployment

The project is designed to be **Kubernetes-ready**. Use the Docker image above and create deployments, services, and ingress resources as needed.

---

## Future Enhancements

* Proxy support for scraping requests
* Multiple session management (cookies, headers)
* Round-robin session selection
* Distributed scraping workers
* Additional social media modules (Instagram, Twitter, TikTok, etc.)

---

## License

MIT License © 2025

---

## Contact

Developed by Er Vijay Raghuwanshi
Email: [ervijayraghuwanshi@gmail.com](mailto:ervijayraghuwanshi@gmail.com)



