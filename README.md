# Go REST API Boilerplate
[![Go Report Card](https://goreportcard.com/badge/github.com/nicobistolfi/go-rest-api)](https://goreportcard.com/report/github.com/nicobistolfi/go-rest-api)
[![Documentation](https://img.shields.io/badge/documentation-yes-blue.svg)](https://go-rest-api.bistol.fi/)
![License](https://img.shields.io/badge/license-MIT-green.svg)
[![Author](https://img.shields.io/badge/author-%40nicobistolfi-blue.svg)](https://github.com/nicobistolfi)

# Introduction

This Go REST API boilerplate establishes a strong foundation for your API development, emphasizing clean architecture, thorough testing, and flexible deployment options. The modular structure ensures maintainability and scalability by promoting a clear separation of concerns, making it easy to modify and extend your API as needed.

Comprehensive testing is essential for a reliable and secure API. The boilerplate includes a full suite of tests, including unit tests, API security tests, service contract tests, and performance benchmarks. With these tests in place, you can confidently deploy your API using the method that best suits your needs. You can find guides on deploying on different platforms in the [Deployments](/docs/deployments) section.

## Getting Started

1. Clone this repository.

2. Navigate to the project root.

3. Run `go mod tidy` to ensure all dependencies are correctly installed.

4. Copy the `.env.example` file to `.env`:
   ```bash
   cp .env.example .env
   ```
5. Open the `.env` file and set the `TOKEN_URL` environment variable to the GitHub API URL:
   ```
   TOKEN_URL=https://api.github.com/user
   ```
   5.1 _If you want to use other providers, you can do so by setting the `TOKEN_URL` environment variable to the provider's API URL and change `token.go` file to use the correct provider._

6. Use the provided Makefile commands for common tasks:
   - `make build`: Build the application
   - `make test`: Run all tests
   - `make run`: Run the application locally

## Documentation

To run the documentation locally:

1. Run the following command:
```bash
make docs
```
2. Open `http://localhost:3001` in your browser

This will start a Docusaurus site with comprehensive project documentation.

## Live Documentation

For the most up-to-date and comprehensive documentation, please visit our [official documentation site](https://go-rest-api.bistol.fi/). This site includes:

- Detailed API references
- In-depth guides on architecture and best practices
- Deployment tutorials for various platforms

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
├── LICENSE
└── README.md
```

<details>
<summary>Directory and File Descriptions</summary>

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
Project documentation, API specifications, and any other relevant documentation.

#### Root Files
- `.gitignore`: Specifies intentionally untracked files to ignore.
- `go.mod` and `go.sum`: Go module files for dependency management.
- `Makefile`: Defines commands for building, testing, and deploying the application.
- `LICENSE`: Contains the MIT License text.
- `README.md`: This file, providing an overview of the project structure.

</details>


## Development Workflow

1. Implement new features or bug fixes in the appropriate packages under `internal/`.
2. Write unit tests in the same package as the code being tested.
3. Write integration tests in the `tests/integration/` directory.
4. Update API documentation in the `docs/` directory as necessary.
5. Use the `scripts/` directory for any automation tasks.
6. Update deployment configurations in `deployments/` if there are infrastructure changes.

## Deployment

This project supports multiple deployment options.

- [Docker](https://go-rest-api.bistol.fi/docs/deployments/docker/)
- [Kubernetes](https://go-rest-api.bistol.fi/docs/deployments/kubernetes/)
- [Serverless](https://go-rest-api.bistol.fi/docs/deployments/serverless/)

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
