# Makefile for Buckly Microservices

# Define Gateways
GATEWAYS = gateway

# Docker image prefix
IMAGE_PREFIX = buckly

# Build output directory
BUILD_OUTPUT_DIR = bin

#Go build flags
GO_BUILD_FLAGS = -ldflags="-s -w"

BUCKLY_API_SWAGGER_URL := http://localhost:8080/swagger/index.html


# Commands

# Generate Swagger documentation for the API
swag-api:
	@echo "Generating Swagger documentation for the API..."
	pushd gateway && swag init --parseDependency --parseInternal --parseDepth 1 -g main.go && popd

clean:
	@echo "Cleaning up build binaries..." 	
	@rm -rf $(BUILD_OUTPUT_DIR)/*

tidy: clean
	@echo "Tidying up Go modules..."
	go mod tidy

build-gateway:
	OUTPUT="$(BUILD_OUTPUT_DIR)/gateway"; \
	echo "Building Gateway..."; \
	env GOOS=linux CGO_ENABLED=0 go build $(GO_BUILD_FLAGS) -o $$OUTPUT ./gateway; \
	echo "Gateway built at $$OUTPUT"

build: tidy swag-api build-gateway

up: build
	@echo "Stopping docker images (if running...)"
	docker compose down
	docker image prune -f
	@echo "Building and starting docker images..."
	docker compose up --build -d
	@echo "Docker images are up and running."
	"$(MAKE)" open-api-url

open-api-url:
	@start "$(BUCKLY_API_SWAGGER_URL)"

