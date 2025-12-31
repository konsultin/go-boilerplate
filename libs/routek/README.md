# routek - YAML-Based Routing Library

A declarative routing library that generates HTTP routes from YAML configuration and maps them to handler methods. Part of the konsultin backend boilerplate.

## Features

- **YAML Configuration**: Define routes in `api-route.yaml`
- **Auto-Mapping**: Maps route handlers to struct methods via reflection
- **Multiple Groups**: Organize routes into logical handler groups
- **Flexible Handlers**: Support multiple handler return signatures
- **Built-in Responder**: Automatic JSON response formatting
- **FastHTTP Integration**: High-performance routing

## Quick Start

**1. Define Routes** (`internal/api-route.yaml`):

```yaml
user:
  route:
    - get: /api/users
      handler: GetAllUsers
    - get: /api/users/:id
      handler: GetUserByID
    - post: /api/users
      handler: CreateUser

auth:
  route:
    - post: /api/auth/login
      handler: Login
```

**2. Create Handlers**:

```go
type UserHandler struct {
    service UserService
}

// Returns (data, error) - automatic success/error response
func (h *UserHandler) GetAllUsers(ctx *fasthttp.RequestCtx) (interface{}, error) {
    return h.service.ListUsers()
}

// Returns error only - automatic error response on failure
func (h *UserHandler) CreateUser(ctx *fasthttp.RequestCtx) error {
    return h.service.CreateUser(ctx.PostBody())
}
```

**3. Register Routes**:

```go
import "github.com/konsultin/project-goes-here/libs/routek"

config := routek.Config{
    RouteFile: "internal/api-route.yaml",
    Handlers: map[string]any{
        "user": userHandler,
        "auth": authHandler,
    },
    Responder: routek.NewResponder(false), // false = production mode
}

router, err := routek.NewRouter(config)
if err != nil {
    log.Fatal(err)
}

fasthttp.ListenAndServe(":8080", router.Handler)
```

## YAML Route Format

```yaml
<group_name>:
  route:
    - <method>: <path>
      handler: <MethodName>
```

**HTTP Methods**: `get`, `post`, `put`, `patch`, `delete`, `head`, `options`

**Path Parameters**:
```yaml
- get: /users/:id
  handler: GetUser
- get: /users/:userId/posts/:postId
  handler: GetUserPost
```

## Handler Signatures

### 1. No Return (Manual Response)
```go
func (h *Handler) Custom(ctx *fasthttp.RequestCtx) {
    ctx.SetStatusCode(200)
    ctx.SetBody([]byte(`{"ok":true}`))
}
```

### 2. Error Only
```go
func (h *Handler) Delete(ctx *fasthttp.RequestCtx) error {
    return h.service.Delete(ctx.UserValue("id").(string))
}
```

### 3. Data + Error (Recommended)
```go
func (h *Handler) Get(ctx *fasthttp.RequestCtx) (interface{}, error) {
    return h.service.Get(ctx.UserValue("id").(string))
}
```

## Configuration

```go
type Config struct {
    RouteFile string         // Path to api-route.yaml (optional)
    Handlers  map[string]any // Handler registry (required)
    Responder *Responder     // Response formatter (optional)
}
```

**Default Route Locations**:
1. `internal/api-route.yaml`
2. `api-route.yaml`
3. `config/api-route.yaml`

## Response Format

**Success** (`(data, nil)`):
```json
{
    "code": "OK",
    "message": "success",
    "data": { }
}
```

**Error**:
```json
{
    "code": "INTERNAL_ERROR",
    "message": "internal server error"
}
```

## Production Example

```go
func main() {
    // Create handlers
    userHandler := handler.NewUserHandler(userService)
    postHandler := handler.NewPostHandler(postService)
    
    // Configure router
    debug := os.Getenv("DEBUG") == "true"
    config := routek.Config{
        Handlers: map[string]any{
            "user": userHandler,
            "post": postHandler,
        },
        Responder: routek.NewResponder(debug),
    }
    
    router, err := routek.NewRouter(config)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
```

## Best Practices

- **Group by Domain**: Organize routes by business domain (user, auth, post)
- **RESTful Paths**: Use `/api/v1/resource` patterns
- **Return Data+Error**: Prefer `(interface{}, error)` for auto-formatting
- **Production Mode**: Set `Responder.Debug = false` in production
- **Error Library**: Use `errk` for structured errors
- **One YAML File**: Keep all routes centralized for easy overview
