# httpk - HTTP Client Library

A production-ready HTTP client library with automatic retry logic, circuit breaker pattern, and comprehensive error handling. Part of the konsultin backend boilerplate.

## Features

- **Retry Logic**: Configurable retry with exponential backoff for failed requests
- **Circuit Breaker**: Prevent cascading failures with circuit breaker pattern
- **Clean API**: Simple methods for GET, POST, PUT, PATCH, DELETE
- **Context Support**: Full support for request context and cancellation
- **Logging Integration**: Optional logger for debugging and monitoring
- **JSON Encoding**: Automatic JSON encoding/decoding

## Quick Start

```go
import (
    "context"
    "github.com/konsultin/project-goes-here/libs/httpk"
)

// Create client with default config
client := httpk.NewClient(httpk.DefaultConfig())

// Make a GET request
ctx := context.Background()
resp, err := client.GET(ctx, "https://api.example.com/users", nil)
if err != nil {
    log.Fatal(err)
}

// Decode JSON response
var users []User
if err := resp.DecodeJSON(&users); err != nil {
    log.Fatal(err)
}
```

## Configuration

### Default Configuration

```go
config := httpk.DefaultConfig()
// Returns:
// - Timeout: 30 seconds
// - Retry: 3 attempts with exponential backoff
// - Circuit Breaker: 5 failures threshold, 60s timeout
```

### Custom Configuration

```go
config := &httpk.Config{
    Timeout: 60 * time.Second,
    
    Retry: &httpk.RetryConfig{
        MaxRetries:     5,
        InitialBackoff: 200 * time.Millisecond,
        MaxBackoff:     10 * time.Second,
        BackoffFactor:  2.0,
        RetryStatuses:  []int{429, 502, 503, 504},
    },
    
    CircuitBreaker: &httpk.CircuitBreakerConfig{
        FailureThreshold: 10,
        SuccessThreshold: 3,
        Timeout:          120 * time.Second,
    },
    
    Logger: myLogger, // Optional logger implementation
}

client := httpk.NewClient(config)
```

## API Reference

### Client Creation

#### `NewClient(cfg *Config) *Client`

Create a new HTTP client with the given configuration. Pass `nil` to use default config.

```go
client := httpk.NewClient(nil) // Use defaults
// or
client := httpk.NewClient(customConfig)
```

### HTTP Methods

#### `GET(ctx context.Context, url string, headers map[string]string) (*Response, error)`

Perform a GET request.

```go
headers := map[string]string{
    "Authorization": "Bearer token",
    "Accept": "application/json",
}

resp, err := client.GET(ctx, "https://api.example.com/data", headers)
```

#### `POST(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error)`

Perform a POST request with JSON body.

```go
payload := map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
}

resp, err := client.POST(ctx, "https://api.example.com/users", payload, nil)
```

#### `PUT(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error)`

Perform a PUT request with JSON body.

```go
update := map[string]interface{}{
    "status": "active",
}

resp, err := client.PUT(ctx, "https://api.example.com/users/123", update, nil)
```

#### `PATCH(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error)`

Perform a PATCH request with JSON body.

```go
resp, err := client.PATCH(ctx, url, partialUpdate, nil)
```

#### `DELETE(ctx context.Context, url string, headers map[string]string) (*Response, error)`

Perform a DELETE request.

```go
resp, err := client.DELETE(ctx, "https://api.example.com/users/123", nil)
```

### Advanced Usage

#### `Do(req *Request) (*Response, error)`

Perform a custom HTTP request.

```go
req := &httpk.Request{
    Method:  "POST",
    URL:     "https://api.example.com/webhook",
    Headers: map[string]string{
        "X-Custom-Header": "value",
    },
    Body: payload,
    Ctx:  ctx,
}

resp, err := client.Do(req)
```

### Response Handling

#### `DecodeJSON(target interface{}) error`

Decode JSON response body into the target struct.

```go
var result struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

if err := resp.DecodeJSON(&result); err != nil {
    log.Fatal(err)
}
```

#### Response Fields

```go
type Response struct {
    StatusCode int           // HTTP status code
    Headers    http.Header   // Response headers
    Body       []byte        // Raw response body
    Duration   time.Duration // Request duration
}
```

## Retry Configuration

### RetryConfig Options

- `MaxRetries` - Maximum number of retry attempts (default: 3)
- `InitialBackoff` - Initial backoff duration (default: 100ms)
- `MaxBackoff` - Maximum backoff duration (default: 5s)
- `BackoffFactor` - Exponential backoff multiplier (default: 2.0)
- `RetryStatuses` - HTTP status codes to retry (default: 429, 502, 503, 504)

### Backoff Calculation

Backoff duration = min(InitialBackoff × BackoffFactor^attempt, MaxBackoff)

Example with defaults:
- Attempt 1: 100ms
- Attempt 2: 200ms
- Attempt 3: 400ms

### Custom Retry Strategy

```go
config := &httpk.Config{
    Retry: &httpk.RetryConfig{
        MaxRetries:     5,
        InitialBackoff: 500 * time.Millisecond,
        MaxBackoff:     30 * time.Second,
        BackoffFactor:  1.5,
        RetryStatuses:  []int{408, 429, 500, 502, 503, 504},
    },
}
```

## Circuit Breaker Pattern

The circuit breaker prevents cascading failures by temporarily blocking requests to a failing service.

### States

1. **CLOSED** - Normal operation, requests pass through
2. **OPEN** - Service is failing, requests are rejected immediately
3. **HALF_OPEN** - Testing if service has recovered

### State Transitions

```
CLOSED --[failures >= threshold]--> OPEN
OPEN --[timeout expired]--> HALF_OPEN
HALF_OPEN --[success >= threshold]--> CLOSED
HALF_OPEN --[any failure]--> OPEN
```

### Configuration

- `FailureThreshold` - Failures needed to open circuit (default: 5)
- `SuccessThreshold` - Successes in half-open to close circuit (default: 2)
- `Timeout` - Time before transitioning from open to half-open (default: 60s)

```go
config := &httpk.Config{
    CircuitBreaker: &httpk.CircuitBreakerConfig{
        FailureThreshold: 10,  // Open after 10 failures
        SuccessThreshold: 3,   // Close after 3 successes
        Timeout:          120 * time.Second,
    },
}
```

### Disable Circuit Breaker

```go
config := &httpk.Config{
    CircuitBreaker: nil, // Disable circuit breaker
}
```

## Logger Integration

Implement the `Logger` interface to integrate with your logging system:

```go
type Logger interface {
    Debugf(format string, args ...interface{})
    Infof(format string, args ...interface{})
    Errorf(format string, args ...interface{})
}

// Use with logk
import "github.com/konsultin/project-goes-here/libs/logk"

config := &httpk.Config{
    Logger: logk.Get(),
}
```

## Production Examples

### API Client Wrapper

```go
type APIClient struct {
    http   *httpk.Client
    apiKey string
}

func NewAPIClient(apiKey string) *APIClient {
    config := &httpk.Config{
        Timeout: 30 * time.Second,
        Retry: &httpk.RetryConfig{
            MaxRetries:     3,
            InitialBackoff: 200 * time.Millisecond,
        },
    }
    
    return &APIClient{
        http:   httpk.NewClient(config),
        apiKey: apiKey,
    }
}

func (c *APIClient) GetUser(ctx context.Context, id string) (*User, error) {
    url := fmt.Sprintf("https://api.example.com/users/%s", id)
    headers := map[string]string{
        "Authorization": "Bearer " + c.apiKey,
        "Accept": "application/json",
    }
    
    resp, err := c.http.GET(ctx, url, headers)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
    
    var user User
    if err := resp.DecodeJSON(&user); err != nil {
        return nil, fmt.Errorf("decode failed: %w", err)
    }
    
    return &user, nil
}
```

### Webhook Sender with Circuit Breaker

```go
type WebhookSender struct {
    client *httpk.Client
}

func NewWebhookSender() *WebhookSender {
    config := &httpk.Config{
        CircuitBreaker: &httpk.CircuitBreakerConfig{
            FailureThreshold: 5,
            SuccessThreshold: 2,
            Timeout:          60 * time.Second,
        },
        Retry: httpk.DefaultRetryConfig(),
    }
    
    return &WebhookSender{
        client: httpk.NewClient(config),
    }
}

func (s *WebhookSender) Send(ctx context.Context, url string, event interface{}) error {
    resp, err := s.client.POST(ctx, url, event, nil)
    if err != nil {
        return fmt.Errorf("webhook failed: %w", err)
    }
    
    if resp.StatusCode >= 400 {
        return fmt.Errorf("webhook returned %d", resp.StatusCode)
    }
    
    return nil
}
```

## Best Practices

- **Use Context**: Always pass context for timeout and cancellation support
- **Handle Status Codes**: Check `resp.StatusCode` before processing responses
- **Configure Timeouts**: Set appropriate timeouts for your use case
- **Enable Circuit Breaker**: Use for external service calls to prevent cascades
- **Retry Sparingly**: Only retry on transient errors (429, 502, 503, 504)
- **Add Logging**: Integrate logger for production debugging
- **Reuse Clients**: Create client instances once and reuse them
