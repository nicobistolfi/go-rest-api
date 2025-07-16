---
sidebar_position: 2
sidebar_label: Taskfile
---

# Using the Taskfile

This document explains how to use the Taskfile provided in the Go REST API Boilerplate project. The Taskfile contains various commands to streamline development, testing, and deployment processes.

## Prerequisites

- Install Go Task: `brew install go-task` (macOS) or see [installation guide](https://taskfile.dev/installation/)
- Ensure you have Go 1.22.5+ installed and properly configured
- Docker should be installed for Docker-related commands

## Common Commands

### Development
```bash
task dev:install  # Install Air for hot reload (first time only)
task dev         # Start development server with hot reload
task build       # Compile the binary with caching
task run         # Build and run in production mode
task clean       # Clean build artifacts
```

### Testing
```bash
task test               # Run all unit tests
task test:coverage      # Run tests with coverage (generates coverage.html)
task test:unit          # Run unit tests (internal/pkg only)
task test:integration   # Run integration tests
task test:performance   # Run performance benchmarks
task test:security      # Run security tests
task test:e2e          # Run end-to-end tests
task test:contract     # Run contract tests
task test:all          # Run all test types
```

### Code Quality
```bash
task fmt            # Format all Go code
task fmt:check      # Check code formatting
task lint           # Run golangci-lint
task vuln:check     # Scan for vulnerabilities
```

### Dependencies
```bash
task mod            # Download and tidy modules
task mod:update     # Update all dependencies
task mod:verify     # Verify dependency integrity
task dep            # Alias for 'task mod'
```

### Docker
```bash
task docker:build   # Build Docker image
task docker:run     # Run Docker container
```

### Documentation
```bash
task docs           # Generate API docs with Swag
task docs:serve     # Start documentation server
```

### Serverless Deployment
```bash
task sls:deploy     # Deploy using Serverless Framework
```

## Quick Start Examples

**Start development:**
```bash
task dev:install  # First time only
task dev
```

**Production build:**
```bash
task mod && task build && task run
```

**Pre-commit checks:**
```bash
task fmt && task lint && task test:all
```

**Docker deployment:**
```bash
task docker:build && task docker:run
```

## Configuration

Customize the Taskfile.yml by modifying these variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `BINARY_NAME` | `api` | Name of compiled binary |
| `BUILD_DIR` | `./build` | Build output directory |
| `MAIN_FILE` | `./cmd/api` | Main entry point |
| `DOCKER_IMAGE` | `my-go-api` | Docker image name |
| `VERSION` | `0.0.1` | Version number |

## Features

- **Build Caching**: Only rebuilds when source files change
- **Preconditions**: Automatically checks for required tools and provides installation instructions
- **Development Mode**: Hot reload with Air, excludes non-Go directories from watch
- **Smart Dependencies**: Tasks automatically run dependencies when needed

Run `task` to see all available commands with descriptions.