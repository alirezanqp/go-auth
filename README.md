# Go OTP Authentication Service

A production-ready OTP-based authentication service built with Go.

## Features

- OTP-based authentication with phone number verification
- Rate limiting (3 requests per phone number in 10 minutes)
- JWT token authentication
- User management with pagination and search
- PostgreSQL database with GORM
- Docker support
- Comprehensive logging
- API versioning
- CI/CD pipeline with GitHub Actions

## Architecture

```
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and migrations
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Data models and DTOs
│   └── services/       # Business logic
└── pkg/utils/          # Reusable utilities
```

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL (if running locally)

## Quick Start

### Using Docker Compose

```bash
git clone <repository-url>
cd go-auth
docker compose up -d
```

### Local Development

```bash
git clone <repository-url>
cd go-auth
go mod download

# Start PostgreSQL
docker run --name postgres-auth \
  -e POSTGRES_DB=go_auth \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 -d postgres:15-alpine

# Run application
go run cmd/server/main.go
```

## Configuration

Environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `go_auth` |
| `JWT_SECRET` | JWT signing secret | `your-super-secret-jwt-key` |
| `PORT` | Server port | `8080` |
| `LOG_LEVEL` | Logging level | `info` |

## API Endpoints

### Authentication

```http
POST /api/v1/auth/send-otp
{
  "phone_number": "+1234567890"
}
```

```http
POST /api/v1/auth/verify-otp
{
  "phone_number": "+1234567890",
  "code": "123456"
}
```

```http
GET /api/v1/auth/profile
Authorization: Bearer <jwt_token>
```

### User Management

```http
GET /api/v1/users?page=1&limit=10&search=+123
Authorization: Bearer <jwt_token>
```

```http
GET /api/v1/users/{user_id}
Authorization: Bearer <jwt_token>
```

```http
GET /api/v1/users/stats
Authorization: Bearer <jwt_token>
```

### System

```http
GET /health
GET /version
GET /api/info
```

## Development Commands

```bash
make run         # Run application
make test        # Run tests with coverage
make lint        # Run linter
make security    # Security scan
make build       # Build application
make docker-up   # Start with Docker
```

## Database

PostgreSQL was chosen for:
- ACID compliance
- JSON support
- Scalability
- UUID support
- Mature Go ecosystem support

## Security

- Rate limiting on OTP requests
- OTP expiration (2 minutes)
- JWT token authentication
- Input validation
- Security event logging
- Vulnerability scanning in CI/CD

## License

MIT License