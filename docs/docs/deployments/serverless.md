---
sidebar_position: 3
---

# Serverless Deployment

This guide explains how to deploy the Go REST API using the Serverless Framework with AWS Lambda.

## Prerequisites

- Node.js and npm installed
- Serverless Framework CLI installed (`npm install -g serverless`)
- AWS CLI configured with your credentials

## Files

The following files in the `/deployments/serverless` directory are used for the serverless deployment:

- `serverless.yml`: Main configuration file for the Serverless Framework
- `Makefile`: Contains commands to build and deploy the application
- `main.go`: The Lambda function entry point

## Deployment Steps

1. Navigate to the serverless deployment directory:
   ```bash
   cd deployments/serverless
   ```

2. Build the Go binary:
   ```bash
   make build
   ```

3. Deploy the application:
   ```bash
   make deploy
   ```

   This command will package your application and deploy it to AWS Lambda using the Serverless Framework.

4. After successful deployment, you'll see output with details about your Lambda function and API Gateway endpoint.

## Configuration

The `serverless.yml` file contains the configuration for your serverless deployment. Key points:

- The service name is set to `go-rest-api`
- It's configured to deploy to the `us-east-1` region by default
- The Lambda function is named `api` and uses the `go1.x` runtime
- An API Gateway is set up to trigger the Lambda function
- The API paths are defined to route requests to the Lambda function

## Makefile Commands

The `Makefile` provides the following commands:

- `make build`: Compiles the Go application for AWS Lambda
- `make deploy`: Deploys the application using Serverless Framework
- `make remove`: Removes the deployed application from AWS

## Customization

To customize the deployment:

1. Modify the `serverless.yml` file to change AWS region, function name, or API Gateway settings.
2. Update the `main.go` file if you need to change the Lambda function's behavior.
3. Adjust the `Makefile` if you need to modify the build or deployment process.

## Accessing the API

After deployment, you can access your API using the URL provided in the Serverless Framework output. It will look something like:

```
https://
