# Minimalistic Go REST API Boilerplate
[![Go Reference](https://pkg.go.dev/badge/github.com/nicobistolfi/go-rest-api.svg)](https://pkg.go.dev/github.com/nicobistolfi/go-rest-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/nicobistolfi/go-rest-api)](https://goreportcard.com/report/github.com/nicobistolfi/go-rest-api)
[![Documentation](https://img.shields.io/badge/documentation-yes-blue.svg)](https://go-rest-api.bistol.fi/)
![License](https://img.shields.io/badge/license-MIT-green.svg)
[![Author](https://img.shields.io/badge/author-%40nicobistolfi-blue.svg)](https://github.com/nicobistolfi)

This Minimalistic Go REST API boilerplate establishes a strong foundation for your API development, emphasizing clean architecture, thorough testing, and flexible deployment options. The modular structure ensures maintainability and scalability by promoting a clear separation of concerns, making it easy to modify and extend your API as needed witouth a steep learning curve on how the API is setup or dependencies.

Comprehensive testing is essential for a reliable and secure API. The boilerplate includes a full suite of tests, including unit tests, API security tests, service contract tests, and performance benchmarks. With these tests in place, you can confidently deploy your API using the method that best suits your needs. You can find guides on deploying on different platforms in the [Deployments](https://go-rest-api.bistol.fi/docs/deployments/) section.

## Principles

- **Simplicity First**: Flat learning curve, easy to understand and modify.
- **Idiomatic Go**: Following the Go proverb "Clear is better than clever."
- **Minimal Dependencies**: Lightweight focused libraries that meet basic needs and can be easily replaced.
- **Modularity**: Flexible and modular, allowing for easy extension without overplanning.
- **Standard Library Preference**: Utilize Go's rich standard library wherever possible, only introducing external dependencies when absolutely necessary.
- **KISS (Keep It Simple, Stupid)**: Avoid premature optimizations and complex abstractions.
- **Rapid Iteration**: Enable fast prototyping and quicker feedback cycles.
- **Error Handling**: Explicit error checks and avoiding panic in regular operations.
- **Concurrency When Appropriate**: Only when they provide clear benefits.
- **Testing Dorito Approach**: Write tests. Not too many. Mostly integration..

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

6. Use the provided Taskfile commands for common tasks:
   - `task build`: Build the application
   - `task test`: Run all tests
   - `task run`: Run the application locally

## Documentation

To run the documentation locally:

1. Run the following command:
```bash
task docs
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
├── Taskfile.yml
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
- `Taskfile.yml`: Defines commands for building, testing, and deploying the application.
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

## Contributing & Pull Requests

### Creating Pull Requests

This project uses a label-based CI/CD system to manage automated testing and quality checks. Here's how to work with pull requests:

#### 🏷️ Required Labels for CI/CD

| Label | Purpose | When to Use |
|-------|---------|-------------|
| `run-checks` | **Required** for any quality checks to run | Add this label to trigger all basic CI/CD checks |
| `e2e-test` | Runs end-to-end tests on Kubernetes | Add when changes affect API endpoints, middleware, or deployment configs |
| `checks-passed` | Automatically added when all checks pass | Don't add manually - the system manages this |

#### 📋 Quality Checks Overview

When you add the `run-checks` label, the following automated checks will run:

1. **🔨 Build & Format**
   - `task build` - Compiles the application
   - `task fmt:check` - Verifies code formatting
   - `task mod:verify` - Validates Go module dependencies

2. **🧪 Comprehensive Testing**
   - `task test:unit` - Unit tests (internal/, pkg/)
   - `task test:integration` - Integration tests
   - `task test:performance` - Performance benchmarks
   - `task test:security` - Security tests
   - `task test:contract` - Contract tests
   - `task test:coverage` - Coverage analysis

3. **🔍 Code Quality**
   - `golangci-lint run` - Linting with 5-minute timeout
   - Excludes: docs, examples, deployments/serverless

4. **🐳 Docker Testing**
   - `task docker:build` - Builds Docker image
   - Container health check on `/api/v1/health`

5. **🔒 Security Scanning**
   - `task vuln:check` - Vulnerability detection with govulncheck

6. **☸️ E2E Testing** (only with `e2e-test` label)
   - Minikube Kubernetes cluster setup
   - `task test:e2e` - Full deployment testing
   - Comprehensive logging on failures

#### 🚀 Pull Request Workflow

1. **Create your feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes and test locally:**
   ```bash
   task fmt           # Format code
   task lint          # Check for issues
   task test:all      # Run all tests
   ```

3. **Push your branch and create a PR:**
   ```bash
   git push origin feature/your-feature-name
   ```

4. **Add appropriate labels:**
   - Always add `run-checks` for basic quality verification
   - Add `e2e-test` if your changes affect:
     - API endpoints or handlers
     - Middleware or authentication
     - Docker or Kubernetes configurations
     - Dependencies that could affect runtime behavior

5. **Monitor the automated quality check comment:**
   - The CI system will post a detailed status comment
   - Green ✅ means all checks passed
   - Red ❌ indicates issues that need attention
   - The `checks-passed` label will be added automatically when all checks succeed

#### 💡 Tips for Successful PRs

- **Local testing first:** Run `task fmt && task lint && task test:all` before pushing
- **Small, focused changes:** Easier to review and test
- **Clear descriptions:** Explain what changed and why
- **Add tests:** New features should include appropriate test coverage
- **Update docs:** Significant changes should update relevant documentation

#### 🔧 Development Commands

```bash
# Quick development cycle
task dev              # Start with hot reload
task test:unit        # Fast unit test feedback
task fmt && task lint # Code quality check

# Pre-commit verification
task fmt && task lint && task test:all

# Local E2E testing (requires Docker)
task docker:build && task docker:run
```

#### 📊 CI/CD Status Messages

The automated system will post status comments showing:
- ✅/❌ Status for each check type
- Direct links to failing job logs
- Task commands that were executed
- Detailed results in expandable sections

For detailed guidelines on code style, commit messages, and review process, please read [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
