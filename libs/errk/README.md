# errk - Structured Error Handling Library

A Go library for structured error handling with namespace, code, metadata, traces, and error wrapping capabilities. Part of the konsultin backend boilerplate.

## Features

- **Structured Errors**: Organize errors with namespace and code for consistent error handling
- **Error Wrapping**: Wrap underlying errors while maintaining error context
- **Tracing**: Automatic trace collection showing where errors originated
- **Metadata**: Attach arbitrary metadata to errors for debugging
- **Error Comparison**: Use Go's standard `errors.Is()` for error type checking

## Quick Start

```go
import "github.com/konsultin/project-goes-here/libs/errk"

// Create a new error
err := errk.NewError("USER_NOT_FOUND", "User with given ID does not exist")

// Create with namespace
err := errk.NewError(
    "VALIDATION_FAILED", 
    "Email format is invalid",
    errk.WithNamespace("auth"),
)

// Add trace information
err = err.Trace()
```

## API Reference

### Creating Errors

#### `NewError(code, message string, options ...SetOptionFn) *Error`

Create a new error instance with the given code and message.

```go
// Basic error
err := errk.NewError("NOT_FOUND", "Resource not found")

// With namespace
err := errk.NewError(
    "DUPLICATE_ENTRY", 
    "Email already exists",
    errk.WithNamespace("database"),
)

// With metadata
err := errk.NewError(
    "RATE_LIMIT_EXCEEDED",
    "Too many requests",
    errk.WithMetadata(map[string]interface{}{
        "limit": 100,
        "window": "1m",
    }),
)
```

### Error Methods

#### `Trace(options ...SetOptionFn) *Error`

Add a trace to the error showing where it occurred. Returns a new error instance.

```go
func validateUser(user User) error {
    if user.Email == "" {
        return errk.NewError("VALIDATION_FAILED", "Email is required").Trace()
    }
    return nil
}
```

#### `Wrap(err error) *Error`

Wrap an existing error as the source error. Useful for adding context to third-party errors.

```go
func loadConfig(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return errk.NewError("CONFIG_LOAD_FAILED", "Failed to load config").Wrap(err)
    }
    return nil
}
```

#### `Copy(options ...SetOptionFn) *Error`

Create a copy of the error with optional modifications.

```go
baseErr := errk.NewError("DB_ERROR", "Database operation failed")

// Copy with different namespace
err := baseErr.Copy(errk.WithNamespace("users"))

// Copy with metadata
err := baseErr.Copy(errk.WithMetadata(map[string]interface{}{
    "query": "SELECT * FROM users",
}))
```

#### `AddMetadata(key string, value interface{}) *Error`

Add metadata to an existing error. Returns a new error instance.

```go
err := errk.NewError("PAYMENT_FAILED", "Payment processing failed")
err = err.AddMetadata("amount", 1500)
err = err.AddMetadata("currency", "USD")
```

### Error Comparison

Use Go's standard `errors.Is()` to check error types:

```go
var ErrNotFound = errk.NewError("NOT_FOUND", "Resource not found", errk.WithNamespace("app"))

func FindUser(id int) error {
    // ... database logic
    return ErrNotFound.Trace()
}

// Check error type
err := FindUser(123)
if errors.Is(err, ErrNotFound) {
    // Handle not found case
}
```

### Getter Methods

- `Code() string` - Get error code
- `Namespace() string` - Get error namespace
- `Message() string` - Get error message
- `Metadata() map[string]interface{}` - Get metadata
- `Traces() []string` - Get trace information

## Production Patterns

### Define Error Constants

```go
package repository

import "github.com/konsultin/project-goes-here/libs/errk"

const namespace = "user.repository"

var (
    ErrUserNotFound = errk.NewError(
        "USER_NOT_FOUND",
        "User does not exist",
        errk.WithNamespace(namespace),
    )
    
    ErrDuplicateEmail = errk.NewError(
        "DUPLICATE_EMAIL",
        "Email already registered",
        errk.WithNamespace(namespace),
    )
)

func (r *Repository) FindByID(id int) (*User, error) {
    user, err := r.db.Get(id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound.Trace()
        }
        return nil, errk.NewError("DB_ERROR", "Database query failed", 
            errk.WithNamespace(namespace)).Wrap(err).Trace()
    }
    return user, nil
}
```

### Error Handling in HTTP Handlers

```go
func (h *Handler) GetUser(ctx *fasthttp.RequestCtx) error {
    id := ctx.UserValue("id").(string)
    
    user, err := h.service.GetUser(id)
    if err != nil {
        // Check error type
        if errors.Is(err, ErrUserNotFound) {
            return h.responder.Error(ctx, 404, "NOT_FOUND", err.Error(), err)
        }
        
        // Log error with traces
        h.logger.Error("Failed to get user", logk.WithError(err))
        
        return h.responder.Error(ctx, 500, "INTERNAL_ERROR", "Operation failed", err)
    }
    
    return h.responder.Success(ctx, 200, "OK", "User retrieved", user)
}
```

## Best Practices

- **Use Namespaces**: Organize errors by domain/package (e.g., `auth`, `database`, `payment`)
- **Consistent Codes**: Use UPPER_SNAKE_CASE for error codes
- **Add Traces**: Always call `.Trace()` when returning errors to capture call stack
- **Wrap External Errors**: Use `.Wrap()` to add context to third-party library errors
- **Define Constants**: Create package-level error constants for reusable errors
- **Include Context**: Use metadata to attach relevant debugging information
