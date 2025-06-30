# Logger Usage Examples

## Basic Usage

```go
import "github.com/conformitea-co/conformitea/internal/logger"

// Simple logging
logger.Info("Server started")
logger.Warn("Configuration missing, using defaults")
logger.Error("Failed to connect to database")

// With fields
logger.Info("User logged in",
    logger.Fields(map[string]interface{}{
        "user_id": "123",
        "ip": "192.168.1.1",
    })...)
```

## In HTTP Handlers

```go
import "github.com/conformitea-co/conformitea/internal/server/middlewares"

func MyHandler(c *gin.Context) {
    // Get logger with request context
    log := middlewares.GetLogger(c)

    // This will include request_id, user_id (if authenticated), etc.
    log.Info("Processing request")

    // Log with additional fields
    log.Info("Database query executed",
        zap.Int("rows_affected", 5),
        zap.Duration("query_time", 150*time.Millisecond),
    )

    // Log errors
    if err != nil {
        log.Error("Failed to process request",
            zap.Error(err),
            zap.String("operation", "user_update"),
        )
    }
}
```

## Configuration Options

### Development (Console Output)

```toml
[logger]
level = "debug"
format = "console"  # Human-readable with colors
output = "stdout"
```

### Production (JSON Output)

```toml
[logger]
level = "info"
format = "json"     # Structured for log aggregation
output = "stdout"   # Or file path like "/var/log/conformitea.log"
```

## Performance Logging

```go
// Log performance metrics
logger.Info("Request completed",
    logger.Performance(
        5,    // db queries
        120,  // db duration ms
        15,   // redis duration ms
        200,  // external API duration ms
    )...)
```

## HTTP Request Logging

The logging middleware automatically logs all HTTP requests with:

- Request ID
- Method and Path
- Status Code
- Response Time
- Client IP
- User Agent
- User ID (if authenticated)
- Errors (if any)

Example output:

```json
{
  "timestamp": "2025-06-30T10:15:30.123Z",
  "level": "info",
  "message": "HTTP request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/auth/login",
  "status": 200,
  "latency_ms": 45,
  "ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "user_id": "user123"
}
```

## Best Practices

1. **Use structured logging**: Always use fields instead of string concatenation

   ```go
   // Good
   logger.Info("User action", zap.String("action", "login"), zap.String("user_id", userID))

   // Bad
   logger.Info(fmt.Sprintf("User %s performed login", userID))
   ```

2. **Use appropriate log levels**:

   - `Debug`: Detailed information for debugging
   - `Info`: General informational messages
   - `Warn`: Warning messages for potentially harmful situations
   - `Error`: Error events that might still allow the app to continue
   - `Fatal`: Severe errors that cause the app to abort

3. **Include context**: Always use the request-scoped logger in handlers

   ```go
   log := middlewares.GetLogger(c)
   ```

4. **Log errors with context**:

   ```go
   log.Error("Operation failed",
       zap.Error(err),
       zap.String("operation", "create_user"),
       zap.Any("input", userInput),
   )
   ```

5. **Use consistent field names**: Stick to conventions like:
   - `user_id` for user identifiers
   - `request_id` for request tracking
   - `error` for error messages
   - `duration_ms` for time measurements
