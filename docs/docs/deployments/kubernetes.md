---
sidebar_position: 2
---

# Kubernetes Deployment

This directory contains Kubernetes manifests for deploying the Go REST API application.

## Files

- `deployment.yaml`: Defines the Deployment for the Go REST API application
- `service.yaml`: Exposes the Deployment as a Service within the cluster
- `configmap.yaml`: Contains configuration data as key-value pairs
- `ingress.yaml`: Configures ingress for external access to the Service

## Local Deployment with Docker Desktop or Minikube

### Prerequisites

- Docker Desktop with Kubernetes enabled, or Minikube installed
- kubectl CLI tool installed

### Steps for Local Deployment

1. Start your local Kubernetes cluster:
   - For Docker Desktop: Enable Kubernetes in Docker Desktop settings
   - For Minikube: Run `minikube start`

2. Build your Docker image locally:
   ```bash
   docker build -t go-rest-api:latest -f ../../deployments/docker/Dockerfile ../..
   ```

3. If using Minikube, load the image into Minikube:
   ```bash
   minikube image load go-rest-api:latest
   ```

4. Update the image in `deployment.yaml`:
   ```yaml
   image: go-rest-api:latest
   ```

5. Apply the Kubernetes manifests:
   ```bash
   kubectl apply -f .
   ```

6. Verify the deployment:
   ```bash
   kubectl get deployments
   kubectl get services
   kubectl get pods
   ```

7. Access the application:
   - The API should now be accessible at `http://localhost:30080`
   - If you're using Docker Desktop with Kubernetes, you might need to use `http://localhost:30080`
   - If you're using Minikube, you might need to run `minikube service go-rest-api-service --url` to get the correct URL

## Deployment Steps (for non-local environments)

1. Ensure you have `kubectl` installed and configured to interact with your Kubernetes cluster.

2. Update the image in `deployment.yaml` to point to your container registry:
   ```yaml
   image: your-registry/go-rest-api:latest
   ```

3. Apply the Kubernetes manifests:
   ```bash
   kubectl apply -f .
   ```

4. Verify the deployment:
   ```bash
   kubectl get deployments
   kubectl get services
   kubectl get pods
   ```

5. If using Ingress, ensure you have an Ingress controller installed in your cluster.

## Scaling

To scale the number of replicas, you can use: