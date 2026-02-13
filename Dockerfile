# -------------------------------------------------
# Build Stage
# -------------------------------------------------
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Install git (required for go modules)
RUN apk add --no-cache git ca-certificates

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Swagger generation (if needed)
# RUN go install github.com/swaggo/swag/cmd/swag@latest

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o osint-scraper cmd/api/main.go

# -------------------------------------------------
# Runtime Stage
# -------------------------------------------------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/osint-scraper /app/osint-scraper

# Expose API port
EXPOSE 8080

# Run as non-root user (distroless default)
USER nonroot:nonroot

# Start the service
ENTRYPOINT ["/app/osint-scraper"]
