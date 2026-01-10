# Konsultin Backend Template

🚀 Production-ready Go backend template. Use this to bootstrap new services with Konsultin's standard libraries and architecture.

## How to Use

1. Click **[Use this template](https://github.com/konsultin/api-template/generate)** above
2. Clone your new repository
3. Initialize the project:

```bash
# 1. Rename module to your project path
make setup-project
# Enter module name (e.g., github.com/konsultin/payment-service)

# 2. Install tools & dependencies
make init

# 3. Start development (hot-reload)
make dev
```

## Features

- **Standardized Architecture** - Clean Architecture (Handler -> Service -> Repository)
- **Modular Core** - Built on `github.com/konsultin/*` libraries
- **Authentication** - Pre-configured JWT, OAuth, and Session management
- **Background Jobs** - NATS-based worker system integrated
- **Dev Experience** - Hot-reload (Air), Docker Compose, Makefiles

## Included Libraries

This template comes pre-wired with:

- [`errk`](https://github.com/konsultin/errk) - Error tracing & wrapping
- [`httpk`](https://github.com/konsultin/httpk) - Resilient HTTP client
- [`logk`](https://github.com/konsultin/logk) - Structured logging
- [`natsk`](https://github.com/konsultin/natsk) - NATS messaging
- [`sqlk`](https://github.com/konsultin/sqlk) - SQL query builder
- [`routek`](https://github.com/konsultin/routek) - YAML routing
- [`timek`](https://github.com/konsultin/timek) - Time utilities

## API Documentation (Swagger 3.0)

This template includes auto-generated API docs using [Swaggo](https://github.com/swaggo/swag).

### 1. Generate Docs
Run this command after modifying API handlers:
```bash
make swagger
```

### 2. View Documentation
Start the server and visit: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### 3. Usage Example
Add annotations above your handler:

```go
// @Summary      Login user
// @Tags         auth
// @Param        request body dto.LoginRequest true "Credentials"
// @Success      200  {object}  dto.Response[dto.LoginResponse]
// @Router       /v1/login [post]
func (s *Server) HandleLogin(ctx *fasthttp.RequestCtx) { ... }
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make swagger` | Generate API docs |
| `make dev` | Start with hot-reload |

## Changelog

### v0.6.0 - Redis Cache
- Added `pkg/redis` helper with connection pool support
- Atomic operations: `SetNX`, `Incr/Decr`
- Configurable timeouts: `DialTimeout`, `ReadTimeout`, `WriteTimeout`
- Added Redis service to `docker-compose.yaml`
- All Docker ports now configurable via `.env`

### v0.5.0 - OpenTelemetry (OTEL)
- Integrated OpenTelemetry (OTEL) for distributed tracing
- Added `pkg/otel` initialization wrapper
- Added HTTP Tracing Middleware for automatic request tracking
- Updated `docker-compose.yaml` with **Jaeger (All-in-One)** for trace visualization
- Enabled by default at `http://localhost:16686`

### v0.4.0 - Cron Sidecar & Swagger
- Integrated Docker application for cron scheduling (Sidecar pattern)
- Added `deployment/cron` with `crontab` and `run.sh` script
- Added unified cron endpoint `GET /v1/cron/:cronType`
- Updated Docker Compose to include `cron` service and `core-api` network alias
- Integrated Swagger (OpenAPI 3.0) for API documentation

### v0.3.0 - Modularization
- Migrated internal libs to independent `github.com/konsultin/*` modules
- Removed local `libs/` directory
- Integrated `routek` for improved routing definition
- Standardized error handling with `errk`


### v0.2.0 - Async Workers
- Integrated NATS for background job processing
- Added worker simulation endpoints
- Standardized publisher/subscriber patterns

### v0.1.0 - Auth System
- Implemented JWT session management
- Added Google OAuth 2.0 support
- Added Anonymous session flow
- Integrated RBAC (Role-Based Access Control)

### v0.0.1 - Genesis
- Initial boilerplate structure
- Docker & Migration setup
- Basic HTTP middleware (CORS, Rate Limit)

## License

MIT License - see [LICENSE](LICENSE) - Created by Kenly Krisaguino