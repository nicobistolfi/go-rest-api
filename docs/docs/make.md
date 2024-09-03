---
sidebar_position: 2
sidebar_label: Makefile
---

# Using the Makefile in Go REST API Boilerplate

This document explains how to use the Makefile provided in the Go REST API Boilerplate project. The Makefile contains various commands to streamline development, testing, and deployment processes.

## Prerequisites

- Make sure you have `make` installed on your system.
- Ensure you have Go installed and properly configured.
- Docker should be installed for Docker-related commands.

## Available Commands

### Building and Running

- `make build`: Compiles the binary.
  ```
  make build
  ```

- `make run`: Builds and runs the binary.
  ```
  make run
  ```

### Cleaning

- `make clean`: Cleans build files and cache.
  ```
  make clean
  ```

### Testing

- `make test`: Runs all unit tests.
  ```
  make test
  ```

- `make test/coverage`: Runs unit tests with coverage.
  ```
  make test/coverage
  ```

- `make test-unit`: Runs only unit tests.
  ```
  make test-unit
  ```

- `make test-integration`: Runs integration tests.
  ```
  make test-integration
  ```

- `make test-performance`: Runs performance tests.
  ```
  make test-performance
  ```

- `make test-security`: Runs security tests.
  ```
  make test-security
  ```

- `make test-e2e`: Runs end-to-end tests.
  ```
  make test-e2e
  ```

- `make test-contract`: Runs contract tests.
  ```
  make test-contract
  ```

- `make test-all`: Runs all types of tests.
  ```
  make test-all
  ```

### Dependencies

- `make dep`: Ensures dependencies are up to date.
  ```
  make dep
  ```

### Code Quality

- `make lint`: Lints the code using golangci-lint.
  ```
  make lint
  ```

### Docker

- `make docker/build`: Builds the Docker image.
  ```
  make docker/build
  ```

- `make docker/run`: Runs the Docker image.
  ```
  make docker/run
  ```

### Documentation

- `make docs`: Starts the documentation server locally.
  ```
  make docs
  ```

### Help

- `make help`: Displays help information about available commands.
  ```
  make help
  ```

## Usage Examples

1. To start development:
   ```
   make dep
   make build
   make run
   ```

2. To run tests before committing:
   ```
   make test-all
   ```

3. To build and run in Docker:
   ```
   make docker/build
   make docker/run
   ```

4. To check code quality:
   ```
   make lint
   ```

5. To view documentation:
   ```
   make docs
   ```

## Customizing the Makefile

You can customize the Makefile by modifying variables at the top:

- `BINARY_NAME`: Change the name of the compiled binary.
- `BUILD_DIR`: Alter the build directory.
- `DOCKER_IMAGE`: Modify the Docker image name.
- `VERSION`: Update the version number.

Remember to run `make help` to see all available commands and their descriptions.