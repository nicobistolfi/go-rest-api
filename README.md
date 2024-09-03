# Go REST API Boilerplate

[![Author](https://img.shields.io/badge/author-%40nicobistolfi-blue.svg)](https://github.com/nicobistolfi)
![License](https://img.shields.io/badge/license-MIT-green.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/nicobistolfi/go-rest-api)](https://goreportcard.com/report/github.com/nicobistolfi/go-rest-api)

This repository provides a structured Go project for building scalable APIs. It emphasizes clean architecture, separation of concerns, and ease of testing and deployment.

## Project Structure

```
go-rest-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── routes.go
│   ├── config/
│   ├── models/
│   ├── repository/
│   └── service/
├── pkg/
├── scripts/
├── tests/
│   ├── integration/
│   └── unit/
├── deployments/
│   ├── docker/
│   ├── kubernetes/
│   └── serverless/
├── docs/
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
├─ LICENSE
└── README.md
```

### Key Components

- `cmd/api/main.go`: Application entry point
- `internal/`: Core application logic
- `pkg/`: Reusable, non-project-specific code
- `tests/`: Comprehensive test suite
- `deployments/`: Deployment configurations for various platforms
- `docs/`: Docusaurus-based documentation

## Functionality

This boilerplate provides:

1. A structured API server setup
2. Configuration management
3. Database integration (via repository layer)
4. Business logic separation (service layer)
5. Middleware support
6. Comprehensive testing framework
7. Multiple deployment options

## Deployment Options

The boilerplate supports multiple deployment strategies:

1. Docker: Containerization for consistent environments
2. Kubernetes: Orchestration for scalable deployments
3. Serverless: Cloud function deployment for serverless architectures

Deployment configurations are located in the `deployments/` directory.

## Documentation

Live documentation for this project can be found at [http://go-rest-api.bistol.fi/](http://go-rest-api.bistol.fi/).

To run the documentation locally:

1. Navigate to the `docs/` directory
2. Run:
   ```
   npm install && npm start
   ```
3. Open `http://localhost:3000` in your browser

This will start a Docusaurus site with comprehensive project documentation.

## Getting Started

1. Clone this repository
2. Run `go mod tidy` to install dependencies
3. Use the provided Makefile commands for common tasks:
   - `make build`: Build the application
   - `make test`: Run all tests
   - `make run`: Run the application locally

For more detailed information, please refer to the full documentation in the `docs/` directory.directory.

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.