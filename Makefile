# Variables
BINARY_NAME=api
BUILD_DIR=build
DOCKER_IMAGE=my-go-api
VERSION?=0.0.1

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## build: Compile the binary.
build:
	@echo "  >  Building binary..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(GOBASE)/cmd/api

## run: Build and run the binary
run: build
	@echo "  >  Running binary..."
	$(BUILD_DIR)/$(BINARY_NAME)

## clean: Clean build files. Runs `go clean` internally.
clean:
	@echo "  >  Cleaning build cache"
	@go clean
	@rm -rf $(BUILD_DIR)

## test: Run unit tests
test:
	@echo "  >  Running unit tests..."
	go test ./...

## test/coverage: Run unit tests with coverage
test/coverage:
	@echo "  >  Running unit tests with coverage..."
	go test ./... -coverprofile=coverage.out

## dep: Get the dependencies
dep:
	@echo "  >  Ensuring dependencies are up to date..."
	@go mod tidy

## lint: Lint the code
lint:
	@echo "  >  Linting code..."
	@golangci-lint run

## docker/build: Build the docker image
docker/build:
	@echo "  >  Building docker image..."
	docker build -t $(DOCKER_IMAGE):$(VERSION) -f deployments/docker/Dockerfile .

## docker/run: Run the docker image
docker/run:
	@echo "  >  Running docker image..."
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(VERSION)

## help: Display this help screen
help:
	@grep -h -E '^[a-zA-Z_/-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build run clean test test/coverage dep lint docker/build docker/run help