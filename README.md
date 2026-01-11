# Konsultin Backend Template

🚀 **Production-ready Go backend template**. Designed for scalability, maintainability, and enterprise-grade observability.

---

## 🏗 Architecture Overview

This project follows **Clean Architecture** principles, ensuring separation of concerns and testability.

### Request Flow
The data flow moves from the outer layer (HTTP) inwards to the core logic and data access.

1. **User/Client** sends an HTTP Request.
2. **Router (routek)** dispatches the request to the appropriate Handler.
3. **Handler (SvcCore)** parses the request (JSON/Query), validates, and invokes `Service`.
4. **Service** executes the business logic (clean of any HTTP/SQL dependencies).
5. **Repository** is called by Service for data operations:
    - Queries **SQL Database**.
    - Publishes events to **NATS**.
    - Calls **3rd Party APIs** (`internal/svc-core/pkg`).
    - Caches data in **Redis**.
    - Stores files in **MinIO/S3**.

### Components

| Layer | Path | Description |
|-------|------|-------------|
| **Handler** | `internal/svc-core/*_server.go` | Entry point for HTTP requests. Parses input, validates headers, calls Service. |
| **Service** | `internal/svc-core/service/` | Pure business logic. Does not know about SQL or HTTP. Testable in isolation. |
| **Repository** | `internal/svc-core/repository/` | Data access layer. Interacts with SQL, Redis, NATS, and external APIs (`internal/svc-core/pkg`). |
| **Model** | `internal/svc-core/model/` | Domain entities and database structs (DTOs for DB). |
| **Pkg** | `internal/svc-core/pkg/` | Shared utilities and 3rd party API integrations. |

---

## 🛠 Project Structure

```bash
├── app/                  # Application entry point (main.go)
├── config/               # Configuration struct and loader
├── deployment/           # Docker and deployment configs
├── dto/                  # Data Transfer Objects (API Request/Response)
├── internal/
│   ├── svc-core/         # Core Microservice logic
│   │   ├── model/        # Database models
│   │   ├── repository/   # Data Access (SQL, Cache, NATS)
│   │   ├── service/      # Business Logic
│   │   └── pkg/          # Service-specific packages & 3rd Party APIs
│   └── middleware/       # HTTP Middlewares (Auth, Logging, etc.)
├── migrations/           # SQL Migration files
├── pkg/                  # Global shared libraries (MinIO, Redis, OTel)
├── Makefile              # Development commands
└── go.mod                # Go Dependencies
```

---

## 🚀 Getting Started

### Prerequisites

- **Go** 1.24+
- **Docker** & **Docker Compose**
- **Make** (Optional, but recommended)
- **Air** (For hot-reload: `go install github.com/air-verse/air@latest`)
- **Migrate** (`go install -tags 'postgres,mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`)

### 1. Project Initialization (First Time Only)
Run this only once when creating a new service from this template.
```bash
# Rename module to your project path (e.g., github.com/my-org/my-service)
make setup-project
```

### 2. Developer Setup (Per Person)
Run this to install necessary tools and dependencies.
```bash
make init
```

### 3. Environment Configuration
Copy the example environment file:
```bash
cp .env.example .env
```

**Key `.env` Variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | `development` | `production` enables strict implementations |
| `DB_DRIVER` | `mysql` | Supports `mysql` or `postgres` |
| `DB_HOST` | `localhost` | Database host |
| `NATS_URL` | `nats://localhost:4222` | NATS Connection URL |
| `REDIS_HOST` | `localhost` | Redis Host |
| `MINIO_ENDPOINT` | `localhost:9000` | MinIO/S3 Endpoint |

### 4. Bootstrap Services
Initialize the environment and services.
```bash
make bs
```

### 5. Start Services
Bring up the Docker services.
```bash
make up
```

### 6. Run Migrations & Start Application
```bash
# Apply Database Migrations
make db-up

# Start API with hot-reload
make dev
```
API will be available at `http://localhost:8080`.

---

## 💻 Development Workflow

### Adding a New Feature

1.  **Define DTO**: Create request/response structs in `dto/`.
2.  **Repository**:
    *   Add SQL queries in `internal/svc-core/sql/`.
    *   Implement methods in `internal/svc-core/repository/`.
3.  **Service**:
    *   Implement business logic in `internal/svc-core/service/`.
    *   Call repository methods.
4.  **Handler**:
    *   Create handler method in `internal/svc-core/xxx_server.go`.
    *   Bind request DTO using `httpk.Bind`.
    *   Call Service method.
5.  **Route**:
    *   Register the handler in `api-route.yaml` or `router setup` in `app/main.go`.

### Database Migrations

```bash
# Create a new migration file
make db-script
# Enter name, e.g., "create_users_table"

# Apply migrations
make db-up

# Rollback last migration
make db-down
```

### API Documentation (Swagger)

1.  **Annotate Handlers**: Add comments above your handler functions.
    ```go
    // @Summary Get User
    // @Tags user
    // @Success 200 {object} dto.UserWrapper
    // @Router /v1/user [get]
    func (s *Server) HandleGetUser(ctx *fasthttp.RequestCtx) { ... }
    ```
2.  **Generate Docs**:
    ```bash
    make swagger
    ```
3.  **View**: Open `http://localhost:8080/swagger/index.html`

---

## 📦 Features

- **Authentication**: JWT & Redis-based Session auth, OAuth (Google).
- **Storage**: MinIO/S3 integration for file uploads.
- **Workers**: Async background jobs using NATS.
- **Observability**: OpenTelemetry (OTel) traces & Structured Logging.
- **Resilience**: Rate limiting, Graceful shutdown, Timeout management.

## 🤝 Contributing

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/amazing-feature`).
3.  Commit your changes (`git commit -m 'Add some amazing feature'`).
4.  Push to the branch (`git push origin feature/amazing-feature`).
5.  Open a Pull Request.

## Changelog

### v0.7.3 - Production Readiness
- Enabled Redis AOF (Append Only File) persistence in `docker-compose.yaml` for data durability
- Clarified MinIO vs S3 usage in documentation

### v0.7.2 - Redis Session Storage
- Migrated User Session storage from PostgreSQL to Redis
- Implemented `DeleteSessionByXid` for Redis compatibility
- Refactored `GetDownloadFileUrl` to `repository/file.go`
- Implemented MinIO S3 for file downloads with configurable expiry
- Updated `pkg/redis` to be context-aware
- Consolidated Database migrations (removed `auth_session`, merged `user_credential`)
- Updated to Go 1.24+

### v0.7.1 - Documentation & Cleanup
- Added comprehensive `README.md` documentation
- Added Architecture Overview diagram and description
- Corrected 3rd party API package location to `internal/svc-core/pkg/`
- Verified and fixed `.env.example` keys

### v0.7.0 - MinIO Storage
- Added `pkg/storage` helper for MinIO/S3 compatible storage
- Upload, Download, Presigned URL helpers
- Auto content-type detection
- Added MinIO service to `docker-compose.yaml`
- Added Redis and MinIO configuration to `config.go`

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

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.