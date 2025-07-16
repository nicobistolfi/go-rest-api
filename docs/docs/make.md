---
sidebar_position: 2
sidebar_label: Taskfile
---

# Using the Taskfile in Go REST API Boilerplate

This document explains how to use the Taskfile provided in the Go REST API Boilerplate project. The Taskfile contains various commands to streamline development, testing, and deployment processes.

## Prerequisites

- Make sure you have `task` installed on your system.
- Ensure you have Go installed and properly configured.
- Docker should be installed for Docker-related commands.

## Available Commands

### Building and Running

- `task build`: Compiles the binary.
  ```
  task build
  ```

- `task run`: Builds and runs the binary.
  ```
  task run
  ```

### Cleaning

- `task clean`: Cleans build files and cache.
  ```
  task clean
  ```

### Testing

- `task test`: Runs all unit tests.
  ```
  task test
  ```

- `task test:coverage`: Runs unit tests with coverage.
  ```
  task test:coverage
  ```

- `task test-unit`: Runs only unit tests.
  ```
  task test-unit
  ```

- `task test-integration`: Runs integration tests.
  ```
  task test-integration
  ```

- `task test-performance`: Runs performance tests.
  ```
  task test-performance
  ```

- `task test-security`: Runs security tests.
  ```
  task test-security
  ```

- `task test-e2e`: Runs end-to-end tests.
  ```
  task test-e2e
  ```

- `task test-contract`: Runs contract tests.
  ```
  task test-contract
  ```

- `task test-all`: Runs all types of tests.
  ```
  task test-all
  ```

### Dependencies

- `task dep`: Ensures dependencies are up to date.
  ```
  task dep
  ```

### Code Quality

- `task lint`: Lints the code using golangci-lint.
  ```
  task lint
  ```

### Docker

- `task docker:build`: Builds the Docker image.
  ```
  task docker:build
  ```

- `task docker:run`: Runs the Docker image.
  ```
  task docker:run
  ```

### Documentation

- `task docs`: Starts the documentation server locally.
  ```
  task docs
  ```

### Help

- `task help`: Displays help information about available commands.
  ```
  task help
  ```

## Usage Examples

1. To start development:
   ```
   task dep
   task build
   task run
   ```

2. To run tests before committing:
   ```
   task test-all
   ```

3. To build and run in Docker:
   ```
   task docker:build
   task docker:run
   ```

4. To check code quality:
   ```
   task lint
   ```

5. To view documentation:
   ```
   task docs
   ```

## Customizing the Taskfile

You can customize the Taskfile by modifying variables at the top:

- `BINARY_NAME`: Change the name of the compiled binary.
- `BUILD_DIR`: Alter the build directory.
- `DOCKER_IMAGE`: Modify the Docker image name.
- `VERSION`: Update the version number.

Remember to run `task help` to see all available commands and their descriptions.