# logk - Logging Library

A flexible logging interface with standard logger implementation, support for multiple log levels, namespaces, and child loggers. Part of the konsultin backend boilerplate.

## Features

- **Interface-Based**: Define custom logger implementations
- **Multiple Log Levels**: Fatal, Error, Warn, Info, Debug, Trace
- **Namespaces**: Organize logs by component or module
- **Child Loggers**: Create isolated loggers with inherited settings
- **Environment Config**: Configure via `LOG_LEVEL` and `LOG_NAMESPACE`
- **Standard Logger**: Built-in implementation using Go's standard logger

## Quick Start

```go
import "github.com/konsultin/project-goes-here/libs/logk"

// Get global logger (auto-initialized from env)
logger := logk.Get()

// Log at different levels
logger.Info("Application started")
logger.Debug("Debug information")
logger.Error("An error occurred")

// Formatted logging
logger.Infof("User %s logged in", username)
logger.Errorf("Failed to connect to %s: %v", host, err)
```

## Environment Variables

- `LOG_LEVEL` - Set log level: `TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL` (default: `INFO`)
- `LOG_NAMESPACE` - Set default namespace/prefix for all logs (default: empty)

Example `.env`:
```bash
LOG_LEVEL=DEBUG
LOG_NAMESPACE=myapp
```

## API Reference

### Logger Interface

```go
type Logger interface {
    // Fatal logs and terminates the application
    Fatal(msg string, options ...logkOption.SetterFunc)
    Fatalf(format string, args ...interface{})
    
    // Error logs error-level messages
    Error(msg string, options ...logkOption.SetterFunc)
    Errorf(format string, args ...interface{})
    
    // Warn logs warning messages
    Warn(msg string, options ...logkOption.SetterFunc)
    Warnf(format string, args ...interface{})
    
    // Info logs informational messages
    Info(msg string, options ...logkOption.SetterFunc)
    Infof(format string, args ...interface{})
    
    // Debug logs debug-level messages
    Debug(msg string, options ...logkOption.SetterFunc)
    Debugf(format string, args ...interface{})
    
    // Trace logs trace-level messages
    Trace(msg string, options ...logkOption.SetterFunc)
    Tracef(format string, args ...interface{})
    
    // NewChild creates a child logger
    NewChild(args ...logkOption.SetterFunc) Logger
}
```

### Global Functions

#### `Get() Logger`

Get the global logger instance. Auto-initializes a standard logger if none is registered.

```go
logger := logk.Get()
logger.Info("Application started")
```

#### `Register(l Logger)`

Register a custom logger implementation.

```go
customLogger := NewCustomLogger()
logk.Register(customLogger)
```

#### `NewChild(args ...logkOption.SetterFunc) Logger`

Create a child logger from the global logger.

```go
dbLogger := logk.NewChild(logkOption.WithNamespace("database"))
dbLogger.Info("Connected to database")
```

#### `Clear()`

Clear the registered logger (mainly for testing).

```go
logk.Clear()
```

## Log Levels

Levels in order of severity (highest to lowest):

1. **FATAL** - Critical errors that terminate the application
2. **ERROR** - Error conditions that need attention
3. **WARN** - Warning messages for potentially harmful situations
4. **INFO** - Informational messages about application flow
5. **DEBUG** - Detailed information for debugging
6. **TRACE** - Very detailed trace information

Setting `LOG_LEVEL=DEBUG` will show DEBUG, INFO, WARN, ERROR, and FATAL logs.

## Using Namespaces

Namespaces help organize logs by component or module:

```go
import (
    "github.com/konsultin/project-goes-here/libs/logk"
    "github.com/konsultin/project-goes-here/libs/logk/option"
)

// Create child logger with namespace
authLogger := logk.NewChild(option.WithNamespace("auth"))
authLogger.Info("User authentication started")
// Output: [auth] User authentication started

dbLogger := logk.NewChild(option.WithNamespace("database"))
dbLogger.Info("Connection established")
// Output: [database] Connection established
```

## Child Loggers

Child loggers inherit settings from parent but can override namespace:

```go
// Parent logger
parentLogger := logk.Get()

// Child with namespace
apiLogger := parentLogger.NewChild(option.WithNamespace("api"))
apiLogger.Info("API server started")

// Grandchild with nested namespace
userAPILogger := apiLogger.NewChild(option.WithNamespace("api.users"))
userAPILogger.Debug("Fetching user list")
```

## Production Examples

### Application Initialization

```go
package main

import (
    "os"
    "github.com/konsultin/project-goes-here/libs/logk"
    "github.com/konsultin/project-goes-here/libs/logk/level"
    "github.com/konsultin/project-goes-here/libs/logk/option"
)

func main() {
    // Logger auto-initialized from environment
    logger := logk.Get()
    
    logger.Info("Application starting...")
    logger.Infof("Running in %s mode", os.Getenv("APP_ENV"))
    
    if err := run(); err != nil {
        logger.Fatalf("Application failed: %v", err)
    }
}
```

### Component-Based Logging

```go
package repository

import (
    "github.com/konsultin/project-goes-here/libs/logk"
    "github.com/konsultin/project-goes-here/libs/logk/option"
)

type UserRepository struct {
    logger logk.Logger
}

func NewUserRepository() *UserRepository {
    return &UserRepository{
        logger: logk.NewChild(option.WithNamespace("repository.user")),
    }
}

func (r *UserRepository) FindByID(id int) (*User, error) {
    r.logger.Debugf("Finding user with ID: %d", id)
    
    user, err := r.db.Get(id)
    if err != nil {
        r.logger.Errorf("Failed to find user %d: %v", id, err)
        return nil, err
    }
    
    r.logger.Infof("Found user: %s", user.Email)
    return user, nil
}
```

### Service Layer Logging

```go
package service

import "github.com/konsultin/project-goes-here/libs/logk"

type AuthService struct {
    logger logk.Logger
    repo   UserRepository
}

func NewAuthService(repo UserRepository) *AuthService {
    return &AuthService{
        logger: logk.NewChild(option.WithNamespace("service.auth")),
        repo:   repo,
    }
}

func (s *AuthService) Login(email, password string) (*Session, error) {
    s.logger.Infof("Login attempt for: %s", email)
    
    user, err := s.repo.FindByEmail(email)
    if err != nil {
        s.logger.Warnf("User not found: %s", email)
        return nil, ErrInvalidCredentials
    }
    
    if !user.ValidatePassword(password) {
        s.logger.Warnf("Invalid password for: %s", email)
        return nil, ErrInvalidCredentials
    }
    
    session := s.createSession(user)
    s.logger.Infof("Login successful: %s", email)
    
    return session, nil
}
```

### HTTP Request Logging

```go
package handler

import (
    "github.com/konsultin/project-goes-here/libs/logk"
    "github.com/valyala/fasthttp"
)

type Handler struct {
    logger logk.Logger
}

func NewHandler() *Handler {
    return &Handler{
        logger: logk.NewChild(option.WithNamespace("handler")),
    }
}

func (h *Handler) GetUser(ctx *fasthttp.RequestCtx) error {
    requestID := string(ctx.Request.Header.Peek("X-Request-ID"))
    
    // Create request-scoped logger
    reqLogger := h.logger.NewChild(option.WithNamespace("handler.getuser"))
    
    reqLogger.Debugf("Request %s: GET /users", requestID)
    
    user, err := h.service.GetUser(ctx)
    if err != nil {
        reqLogger.Errorf("Request %s failed: %v", requestID, err)
        return err
    }
    
    reqLogger.Infof("Request %s completed successfully", requestID)
    return nil
}
```

### Structured Logging with Options

```go
import (
    "github.com/konsultin/project-goes-here/libs/logk"
    "github.com/konsultin/project-goes-here/libs/logk/option"
)

func ProcessPayment(amount float64, currency string) {
    logger := logk.Get()
    
    logger.Info("Processing payment", 
        option.WithMetadata(map[string]interface{}{
            "amount": amount,
            "currency": currency,
        }),
    )
}
```

## Custom Logger Implementation

You can implement your own logger (e.g., for JSON logging, third-party services):

```go
type CustomLogger struct {
    namespace string
    level     level.Level
}

func (l *CustomLogger) Info(msg string, options ...option.SetterFunc) {
    // Your implementation
}

func (l *CustomLogger) Infof(format string, args ...interface{}) {
    // Your implementation
}

// ... implement other methods

// Register your logger
logk.Register(NewCustomLogger())
```

## Best Practices

- **Use Namespaces**: Create child loggers with namespaces for each component
- **Appropriate Levels**: Use correct log levels (don't log everything as Error)
- **Structured Logging**: Include context like user IDs, request IDs in messages
- **Avoid Sensitive Data**: Never log passwords, tokens, or PII
- **Child Loggers**: Use `NewChild()` to create component-specific loggers
- **Production Levels**: Set `LOG_LEVEL=INFO` or `WARN` in production
- **Development Levels**: Use `DEBUG` or `TRACE` for local development
- **Log Errors**: Always log errors with context before returning them
