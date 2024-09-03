# Environment Variables

The authentication system uses environment variables for configuration. These variables allow for flexible deployment across different environments.

## Key Environment Variables

1. `TOKEN_URL`
   - Purpose: Specifies the URL of the external authentication service used for token validation.
   - Required: Yes
   - Example: `https://auth.example.com/validate`

2. `TOKEN_CACHE_EXPIRY`
   - Purpose: Sets the expiration time for cached tokens.
   - Usage: Parsed in the `init()` function of `token.go`
   - Required: No (defaults to 5 minutes if not set)
   - Example: `15m` for 15 minutes, `1h` for 1 hour

## Setting Environment Variables

### In Development

Use a `.env` file in the project root:

```
TOKEN_URL=https://auth.example.com/validate
TOKEN_CACHE_EXPIRY=10m
```

### In Production

Set environment variables according to your deployment platform:

- Docker: Use the `-e` flag or `environment` section in docker-compose.yml
- Kubernetes: Use `ConfigMap` and `Secret` resources
- Cloud Platforms: Use the platform-specific method for setting environment variables