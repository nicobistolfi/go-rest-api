---
title: Token Caching
---

# Token Caching

The token caching mechanism is implemented within the Token Middleware to improve performance by reducing the number of external API calls for token validation.

## Key Features

- In-memory cache for validated tokens and user profiles
- Configurable cache expiry time
- Thread-safe cache operations using a read-write mutex

## Cache Structure

```go
type cacheEntry struct {
    profile Profile
    expiry  time.Time
}

var (
    tokenCache  = make(map[string]cacheEntry)
    cacheMutex  sync.RWMutex
    cacheExpiry time.Duration
)
```

## Cache Operations

1. **Cache Check**: Before validating a token, the middleware checks if a valid cache entry exists.
2. **Cache Hit**: If a non-expired entry is found, it's used directly, skipping external validation.
3. **Cache Miss**: If no valid entry is found, the token is validated externally, and the result is cached.
4. **Cache Update**: After successful validation, the token and profile are cached with an expiry time.

## Cache Expiry Configuration

The cache expiry time can be configured using the `TOKEN_CACHE_EXPIRY` environment variable. If not set, it defaults to 5 minutes.

## Cache Usage in Token Verification

The caching mechanism significantly reduces the load on the external authentication service and improves response times for subsequent requests with the same token. If the cache is not hit, the token is validated externally and the result is cached.

Header `X-Token-Cache` is set to `HIT` if the cache is hit, and `MISS` if the cache is not hit


