# OAuth 2.0 with OpenID Connect (OIDC) Authentication

This document describes the OAuth 2.0 with OpenID Connect (OIDC) authentication mechanism implemented in our API.

## Overview

OAuth 2.0 with OIDC is used for secure user authentication and authorization. It allows users to log in using their accounts from an OIDC provider (e.g., Google, Auth0) without sharing their credentials with our application.

## Configuration

The following environment variables need to be set:

- `OIDC_ISSUER`: The URL of the OIDC provider
- `OAUTH_CLIENT_ID`: The client ID provided by the OIDC provider
- `OAUTH_CLIENT_SECRET`: The client secret provided by the OIDC provider
- `OAUTH_REDIRECT_URL`: The URL to redirect to after successful authentication

## Usage

1. The client initiates the OAuth flow by redirecting the user to the OIDC provider's authorization endpoint.
2. After successful authentication, the OIDC provider redirects back to our application with an authorization code.
3. Our application exchanges the authorization code for an ID token and access token.
4. The access token is then used to authenticate API requests.

To authenticate API requests, include the access token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

## Implementation Details

- The OAuth middleware is implemented in `internal/api/middleware/oauth.go`.
- Token validation and user info retrieval are handled in `pkg/auth/oauth.go`.
- The middleware validates the token and sets the user information in the request context.

## Security Considerations

- Always use HTTPS to protect token transmission.
- Implement proper token storage and management on the