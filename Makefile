.PHONY: all build clean test coverage lint proto migrate docker-build docker-run help

# Variables
BINARY_NAME=hcp-server
DOCKER_IMAGE=fffeng99999/hcp-server
VERSION?=1.0.0
BUILD_DIR=bin

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)
GOFILES=$(wildcard *.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

all: build

## build: Build the binary
build:
	@echo "  >  Building binary..."
	@go build -o $(GOBIN)/$(BINARY_NAME) cmd/server/main.go

## clean: Clean build files
clean:
	@echo "  >  Cleaning build cache"
	@go clean
	@rm -rf $(BUILD_DIR)

## test: Run unit tests
test:
	@echo "  >  Running unit tests..."
	@go test -v -short ./tests/unit/...

## integration-test: Run integration tests
integration-test:
	@echo "  >  Running integration tests..."
	@go test -v ./tests/integration/...

## coverage: Run tests with coverage
coverage:
	@echo "  >  Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

## deps: Download dependencies
deps:
	@echo "  >  Downloading dependencies..."
	@go mod download
	@go mod tidy

## proto: Generate protobuf files
proto:
	@echo "  >  Generating protobuf files..."
	@./scripts/generate_proto.sh

## migrate: Run database migrations
migrate:
	@echo "  >  Running database migrations..."
	@go run cmd/server/main.go migrate

## run: Run the server
run:
	@echo "  >  Running server..."
	@go run cmd/server/main.go

## docker-build: Build docker image
docker-build:
	@echo "  >  Building docker image..."
	@docker build -t $(DOCKER_IMAGE):$(VERSION) -f deployments/Dockerfile .

## docker-run: Run docker container
docker-run:
	@echo "  >  Running docker container..."
	@docker run -p 8081:8081 $(DOCKER_IMAGE):$(VERSION)

## help: Show help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
