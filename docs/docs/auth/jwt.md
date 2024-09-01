# JSON Web Token (JWT) Authentication

This document describes the JWT authentication mechanism implemented in our API.

## Overview

JWT authentication uses a signed token to securely transmit information between parties as a JSON object. This information can be verified and trusted because it is digitally signed.

## Configuration

The following environment variables need to be set:

- `JWT_SECRET`: A secret key used to sign and verify JWTs
- `JWT_EXPIRATION_MINUTES`: The expiration time for JWTs in minutes (default: 60)

## Usage

1. The client obtains a JWT by authenticating with valid credentials (e.g., username and password).
2. The server generates and returns a JWT.
3. The client includes this token in the `Authorization` header for subsequent API requests: