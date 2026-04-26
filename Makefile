# Makefile for Buckly Microservices

# Define Gateways
GATEWAYS = gateway

# Define Services
SERVICES = database-service auth-service

# Docker image prefix
IMAGE_PREFIX = buckly

# Build output directory
BUILD_OUTPUT_DIR = bin

#Go build flags
GO_BUILD_FLAGS = -ldflags="-s -w"

BUCKLY_API_SWAGGER_URL := http://localhost:8080/swagger/index.html


# Commands

sqlc-gen:
	@cd services/database-service && sqlc generate

generate-proto:
	@echo "Generating gRPC code from proto files..."
	@protoc --go_out=. --go-grpc_out=. proto/*/*.proto
	@protoc --go_out=. --go-grpc_out=. proto/*/*/*.proto


# Generate Swagger documentation for the API
swag-api:
	@echo "Generating Swagger documentation for the API..."
	pushd gateway && swag init --parseDependency --parseInternal --parseDepth 1 -g main.go && popd

# Generate certificates for gRPC services
mkcert-services:
	@for service in ${SERVICES}; do \
		echo "Generating certificates for $$service..."; \
		mkcert -cert-file ./certs/$$service.pem -key-file ./certs/$$service-key.pem $$service; \
	done

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

build-services:
	@for service in $(SERVICES); do \
  		echo "Building $$service..."; \
  		cd services/$$service && env GOOS=linux CGO_ENABLED=0 go build $(GO_BUILD_FLAGS) -o ../../bin/$$service . || exit 1; \
  		cd - >/dev/null; \
	done

build: tidy swag-api build-gateway build-services

up: generate-proto build
	@echo "Stopping docker images (if running...)"
	docker compose down
	docker image prune -f
	@echo "Building and starting docker images..."
	docker compose up --build -d
	@echo "Docker images are up and running."
	"$(MAKE)" open-api-url

open-api-url:
	@start "$(BUCKLY_API_SWAGGER_URL)"

# Create a new migration file
# Usage: make create-migration NAME=create_users_table
create-migration:
	@if [ -z "$(NAME)" ]; then \
			echo "Error: NAME is required"; \
			echo "Usage: make create-migration NAME=migration_name"; \
			exit 1; \
	fi
	@goose -dir services/database-service/db/migrations create $(NAME) sql
	@echo "Migration created successfully!"

migrate-up:
	@echo "Running DB migrations..."
	@cd services/database-service/db/migrations && goose postgres "postgres://postgres:External000@localhost:5433/buckly-dev-db?sslmode=disable" up

migrate-down:
	@cd services/database-service/db/migrations && goose postgres "postgres://postgres:External000@localhost:5433/buckly-dev-db?sslmode=disable" down

migrate-status:
	@cd services/database-service/db/migrations && goose postgres "postgres://postgres:External000@localhost:5433/buckly-dev-db?sslmode=disable" status