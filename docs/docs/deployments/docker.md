---
sidebar_position: 2
---

# Docker Deployment

This guide explains how to deploy the Go REST API using Docker.

## Prerequisites

- Docker installed on your machine
- Basic knowledge of Docker commands

## Files

The following file in the `/deployments/docker` directory is used for the Docker deployment:

- `Dockerfile`: Contains instructions for building the Docker image

## Deployment Steps

1. Navigate to the project root directory:
   ```bash
   cd go-rest-api
   ```

2. Build the Docker image:
   ```bash
   docker build -t go-rest-api:latest -f deployments/docker/Dockerfile .
   ```

3. Run the Docker container:
   ```bash
   docker run -p 8080:8080 go-rest-api:latest
   ```

   This command will start the container and map port 8080 from the container to port 8080 on your host machine.

4. The API should now be accessible at `http://localhost:8080`

## Dockerfile Explanation

The `Dockerfile` uses a multi-stage build process:

1. Build stage:
   - Uses the official Go image as the base
   - Sets the working directory
   - Copies the Go module files and downloads dependencies
   - Copies the source code and builds the application

2. Final stage:
   - Uses a minimal Alpine Linux image
   - Copies the built binary from the build stage
   - Sets the entry point to run the application

This approach results in a smaller final image that contains only the necessary components to run the application.

## Customization

To customize the Docker deployment:

1. Modify the `Dockerfile` if you need to change the build process or add additional dependencies.
2. Adjust the exposed port in both the `Dockerfile` and the `docker run` command if you want to use a different port.

## Additional Docker Commands

- To stop the running container:
  ```bash
  docker stop $(docker ps -q --filter ancestor=go-rest-api:latest)
  ```

- To remove the Docker image:
  ```bash
  docker rmi go-rest-api:latest
  ```

- To view logs of the running container:
  ```bash
  docker logs $(docker ps -q --filter ancestor=go-rest-api:latest)
  ```

## Docker Compose (Optional)

If you want to use Docker Compose for easier management, especially when integrating with other services, you can create a `docker-compose.yml` file in the project root:

