# Go API Project Structure

This repository contains a structured Go project for developing a robust and scalable API. The project is organized to promote clean architecture, separation of concerns, and ease of testing and deployment.

## Project Structure

```
go-boilerplate/
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

### Directory and File Descriptions

#### `cmd/`
Contains the main applications for this project. The `api/` subdirectory is where the main.go file for starting the API server resides.

- `api/main.go`: Entry point of the application. Initializes and starts the API server.

#### `internal/`
Houses packages that are specific to this project and not intended for external use.

- `api/`: Contains API-specific code.
  - `handlers/`: Request handlers for each API endpoint.
  - `middleware/`: Custom middleware functions.
  - `routes.go`: Defines API routes and links them to handlers.
- `config/`: Configuration management for the application.
- `models/`: Data models and DTOs (Data Transfer Objects).
- `repository/`: Data access layer, interfacing with the database.
- `service/`: Business logic layer, implementing core functionality.

#### `pkg/`
Shared packages that could potentially be used by external projects. Place reusable, non-project-specific code here.

#### `scripts/`
Utility scripts for development, CI/CD, database migrations, etc.

#### `tests/`
Contains test files separated into integration and unit tests.

- `integration/`: API-level and end-to-end tests.
- `unit/`: Unit tests for individual functions and methods.

#### `deployments/`
Configuration files and scripts for deploying the application.

- `docker/`: Dockerfile and related configurations for containerization.
- `kubernetes/`: Kubernetes manifests for orchestration.
- `serverless/`: Serverless configuration files for cloud function deployment.

#### `docs/`
Project documentation, API specifications, and any other relevant documentation. This directory contains a Docusaurus project for easy-to-navigate and visually appealing documentation.

To run and view the documentation locally:

1. Navigate to the `docs/` directory:
   ```
   cd docs
   ```

2. Install the necessary dependencies:
   ```
   npm install
   ```

3. Start the Docusaurus development server:
   ```
   npm run start
   ```

4. Open your web browser and visit `http://localhost:3000` to view the documentation.

The documentation includes:
- API specifications
- Getting started guide
- Architecture overview
- Deployment instructions
- Contributing guidelines

To build the documentation for production:
```
npm run build
```

This will generate static content in the `build` directory, which can be served using any static content hosting service.

#### Root Files
- `.gitignore`: Specifies intentionally untracked files to ignore.
- `go.mod` and `go.sum`: Go module files for dependency management.
- `Makefile`: Defines commands for building, testing, and deploying the application.
- `LICENSE`: Contains the MIT License text.
- `README.md`: This file, providing an overview of the project structure.

## Getting Started

1. Clone this repository.
2. Navigate to the project root.
3. Run `go mod tidy` to ensure all dependencies are correctly installed.
4. Use the provided Makefile commands for common tasks:
   - `make build`: Build the application
   - `make test`: Run all tests
   - `make run`: Run the application locally

## Development Workflow

1. Implement new features or bug fixes in the appropriate packages under `internal/`.
2. Write unit tests in the same package as the code being tested.
3. Write integration tests in the `tests/integration/` directory.
4. Update API documentation in the `docs/` directory as necessary.
5. Use the `scripts/` directory for any automation tasks.
6. Update deployment configurations in `deployments/` if there are infrastructure changes.

## Deployment

This project supports multiple deployment options:

### Docker

Refer to the `deployments/docker/` directory for Docker configurations. To build and run the Docker container:

1. Build the Docker image:
   ```
   docker build -t go-rest-api -f deployments/docker/Dockerfile .
   ```
2. Run the container:
   ```
   docker run -p 8080:8080 go-rest-api
   ```

### Kubernetes

Kubernetes manifests are available in the `deployments/kubernetes/` directory. To deploy to a Kubernetes cluster:

1. Apply the manifests:
   ```
   kubectl apply -f deployments/kubernetes/
   ```

### Serverless

For serverless deployment, we use the Serverless Framework. Configuration files are located in the `deployments/serverless/` directory.

1. Install the Serverless Framework:
   ```
   npm install -g serverless
   ```
2. Deploy the application:
   ```
   cd deployments/serverless
   serverless deploy
   ```

Ensure to update these configurations as the application evolves. For more detailed deployment instructions, refer to the respective README files in each deployment directory.

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Testing

This project employs a comprehensive testing strategy to ensure reliability, performance, and security. For detailed information on our testing approach, please refer to the following documentation:

- [Testing Strategy](docs/docs/tests/testing-strategy.md)
- [Unit Testing](docs/docs/tests/unit-testing.md)
- [Integration Testing](docs/docs/tests/integration-testing.md)
- [Performance Testing](docs/docs/tests/performance-testing.md)
- [Security Testing](docs/docs/tests/security-testing.md)
- [End-to-End Testing](docs/docs/tests/e2e-testing.md)

To run the tests, you can use the following make commands:

- Run unit tests: `make test-unit`
- Run integration tests: `make test-integration`
- Run performance tests: `make test-performance`
- Run security tests: `make test-security`
- Run end-to-end tests: `make test-e2e`
- Run all tests: `make test-all`

For more information on these commands, refer to the Makefile in the root directory of the project.
