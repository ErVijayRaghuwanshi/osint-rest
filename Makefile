# -------------------------------------------------
# Project Variables
# -------------------------------------------------
APP_NAME := osint-scraper

# Read version safely; fallback to 'latest' if file missing or empty
VERSION  := $(shell [ -f version ] && cat version || echo latest)

MAIN_FILE := cmd/api/main.go
SWAG_DIR  := docs
BIN_DIR   := bin

# -------------------------------------------------
# Help
# -------------------------------------------------
help:
	@echo "Available commands:"
	@echo "  deps           - Install dependencies"
	@echo "  tidy           - Tidy Go modules"
	@echo "  swagger        - Generate Swagger docs"
	@echo "  run            - Run application"
	@echo "  build          - Build local binary"
	@echo "  build-linux    - Build Linux binary (Docker/K8s)"
	@echo "  test           - Run tests"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  k8s-apply      - Apply Kubernetes manifests"
	@echo "  k8s-delete     - Delete Kubernetes resources"
	@echo "  clean          - Clean build artifacts"
	@echo "  all            - deps + tidy + swagger + build"

# -------------------------------------------------
# Go Commands
# -------------------------------------------------

deps:
	go get github.com/gin-gonic/gin
	go get github.com/swaggo/gin-swagger
	go get github.com/swaggo/files
	go get github.com/rs/zerolog
	go install github.com/swaggo/swag/cmd/swag@latest

tidy:
	go mod tidy

swagger:
	@echo "Generating Swagger documentation..."
	rm -rf $(SWAG_DIR)
	swag init -g $(MAIN_FILE) -o $(SWAG_DIR)

run: swagger
	@echo "Starting API server..."
	go run $(MAIN_FILE)

build:
	@echo "Building binary..."
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_FILE)

build-linux:
	@echo "Building Linux binary..."
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_FILE)

test:
	go test ./... -v

# -------------------------------------------------
# Minikube / Kubernetes Docker
# -------------------------------------------------

## Point Docker CLI to Minikube Docker daemon
minikube-docker-env:
	eval $$(minikube docker-env)

## Build Docker image inside Minikube
docker-build-k8s:
	@echo "Building Docker image inside Minikube..."
	eval $$(minikube docker-env) && \
	docker build -t $(APP_NAME):$(VERSION) .

## Deploy to Kubernetes (Minikube)
k8s-deploy: docker-build-k8s
	kubectl apply -f k8s/

## Restart deployment (pick up new image)
k8s-restart:
	kubectl rollout restart deployment/$(APP_NAME) -n osint

# -------------------------------------------------
# Kubernetes
# -------------------------------------------------

k8s-apply:
	kubectl apply -f k8s/

k8s-delete:
	kubectl delete -f k8s/

# -------------------------------------------------
# Utility
# -------------------------------------------------

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(SWAG_DIR)

all: deps tidy swagger build
