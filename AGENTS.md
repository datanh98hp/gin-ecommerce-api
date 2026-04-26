# AGENTS.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

This is a production-ready e-commerce API built with Go and the Gin framework. The architecture follows clean architecture principles with clear separation between API handlers, middleware, models, database, and utilities.

## Key Architecture Decisions

### Layered Structure
- **cmd/**: Application entry point with server initialization and graceful shutdown
- **internal/api/**: HTTP layer (handlers, middleware, routes)
- **internal/models/**: GORM models with validation tags and relationships
- **internal/database/**: Database connection, migrations via GORM AutoMigrate
- **internal/config/**: Environment-based configuration with sensible defaults
- **pkg/utils/**: Reusable utilities (JWT, password hashing, response formatting)

### Authentication Flow
JWT tokens are generated on login/register and validated in `middleware.AuthMiddleware()`. User context (user_id, user_email, user_role) is stored in gin.Context and accessed via `c.Get()`. Admin routes use both `AuthMiddleware` and `AdminMiddleware`.

### Database Relationships
- User → Cart (1:1), User → Orders (1:many), User → Posts (1:many)
- Cart → CartItems (1:many), CartItem → Product (many:1)
- Order → OrderItems (1:many), OrderItem → Product (many:1)
- All relationships use GORM's `Preload()` for eager loading when needed

### Transaction Management
Critical operations (order creation, order cancellation) use `database.DB.Transaction()` to ensure atomicity. Order creation reduces product stock and clears cart; order cancellation restores stock.

## Common Commands

### Running the Application
```bash
# Without Docker
go run cmd/main.go

# With Docker (development)
docker-compose -f docker-compose.dev.yml up --build

# With Docker (staging)
docker-compose -f docker-compose.staging.yml up --build

# With Docker (production)
docker-compose -f docker-compose.prod.yml up --build
```

### Development
```bash
# Install/update dependencies
go mod download
go mod tidy

# Build binary
go build -o main cmd/main.go

# Build for production (static binary)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Format code
go fmt ./...

# Run tests (when added)
go test ./...
go test -v ./internal/...
```

### Docker Management
```bash
# View logs
docker-compose -f docker-compose.dev.yml logs -f app

# Stop and remove containers
docker-compose -f docker-compose.dev.yml down

# Fresh start (removes volumes)
docker-compose -f docker-compose.dev.yml down -v

# Rebuild without cache
docker-compose -f docker-compose.dev.yml build --no-cache
```

## Configuration

Environment variables are loaded via `internal/config/config.go` with fallback defaults. See `.env.example` for all available options. Critical settings:
- `ENV`: Controls Gin mode (development/staging/production) and logging verbosity
- `JWT_SECRET`: Must be changed in production
- `DB_*`: PostgreSQL connection parameters

## Adding New Features

### New API Endpoint
1. Define request/response models in `internal/models/`
2. Create handler method in appropriate handler file under `internal/api/handlers/`
3. Register route in `internal/api/routes/routes.go` (public or protected group)
4. Use `utils.SuccessResponse()` and `utils.ErrorResponse()` for consistent responses

### New Database Model
1. Define struct in `internal/models/` with GORM tags
2. Add to `database.Migrate()` in `internal/database/database.go`
3. Use soft deletes (`gorm.DeletedAt`) for user-facing data
4. Include proper indexes via tags (e.g., `gorm:"uniqueIndex"`)

### New Middleware
1. Create in `internal/api/middleware/`
2. Follow pattern: return `gin.HandlerFunc` that calls `c.Next()` or `c.Abort()`
3. Apply globally in `routes.SetupRoutes()` or to specific route groups

## Important Patterns

### Error Handling
Handlers use early returns with error responses. Never expose raw database errors to clients - return generic messages via `utils.ErrorResponse()`.

### User Authorization
Ownership checks compare `userID.(uint)` from context against resource's UserID. Admin role bypasses ownership checks. Example in `handlers/post.go` Update/Delete methods.

### Stock Management
Product stock is decremented during order creation and restored on cancellation using `gorm.Expr("stock - ?", quantity)` to avoid race conditions.

## Testing Strategy

The application is structured to support testing:
- Handlers can be tested by mocking database.DB
- Middleware can be tested with gin test contexts
- Utils are pure functions easily unit tested
- Integration tests should use a test database

## Database Access

Direct access to `database.DB` is used throughout handlers. For testing or more complex applications, consider the repository pattern (directories exist but are unused).

## API Versioning

All routes are under `/api/v1` prefix. When adding v2, create new route group and handlers while maintaining v1 for backward compatibility.
