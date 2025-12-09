# -------------------------------------------------
# Project Variables
# -------------------------------------------------
APP_NAME=osint-scraper
MAIN_FILE=cmd/api/main.go
SWAG_DIR=docs

# -------------------------------------------------
# Go Commands
# -------------------------------------------------

## Initialize Go module
init:
	@echo "Initializing Go module..."
	go mod init $(APP_NAME)

## Tidy modules
tidy:
	go mod tidy

## Install dependencies
deps:
	go get github.com/gin-gonic/gin
	go get github.com/swaggo/gin-swagger
	go get github.com/swaggo/files
	go get github.com/rs/zerolog
	go get -u github.com/swaggo/swag
# 	go install github.com/swaggo/swag/cmd/swag@latest

## Generate Swagger docs
swagger:
	@echo "Generating Swagger docs..."
	swag init -g $(MAIN_FILE) -o $(SWAG_DIR)

## Run the application
run: swagger
	@echo "Running application..."
	go run $(MAIN_FILE)

## Build binary for local dev
build:
	@echo "Building binary..."
	go build -o bin/$(APP_NAME) $(MAIN_FILE)

## Build static binary for Linux (for Docker/K8s)
build-linux:
	@echo "Building Linux binary..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(APP_NAME) $(MAIN_FILE)

# -------------------------------------------------
# Docker Commands
# -------------------------------------------------

## Build Docker image
docker-build:
	docker build -t $(APP_NAME):latest .

## Run Docker container
docker-run:
	docker run -p 8080:8080 $(APP_NAME):latest

# -------------------------------------------------
# Kubernetes Commands (placeholder)
# -------------------------------------------------

## Apply K8s manifests
k8s-apply:
	kubectl apply -f k8s/

## Delete K8s resources
k8s-delete:
	kubectl delete -f k8s/

# -------------------------------------------------
# Utility
# -------------------------------------------------

## Clean build files
clean:
	rm -rf bin/
	rm -rf $(SWAG_DIR)

## Full build workflow
all: deps tidy swagger build
