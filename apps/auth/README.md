# Authentication Service

Authentication and authorization service for Smart Document platform.

## Overview

The Auth service handles user authentication, authorization, and token management. It provides both HTTP API and gRPC endpoints for other services.

## Architecture

- **API Component**: HTTP endpoints for user-facing authentication
- **RPC Component**: gRPC service for internal service-to-service communication

## Features

- User registration and login
- JWT token generation and validation
- Password hashing and verification
- User profile management
- Token refresh mechanism
- Role-based access control (RBAC)

## Components

### API Service (`/api`)
HTTP REST API for client applications.

**Endpoints:**
- `POST /login` - User login
- `POST /register` - User registration
- `POST /refresh` - Token refresh
- `GET /profile` - User profile
- `PUT /profile` - Update profile

### RPC Service (`/rpc`)
gRPC service for internal communication.

**Methods:**
- `Login(LoginRequest) -> LoginResponse`
- `VerifyToken(TokenRequest) -> VerifyResponse`
- `GetUser(UserRequest) -> UserResponse`

## Configuration

### API Configuration (`etc/auth-api.yaml`)
- Server port and host
- Database connection
- JWT settings
- Rate limiting

### RPC Configuration (`etc/auth-rpc.yaml`)
- gRPC server settings
- Service discovery
- Database connection

## Environment Variables

- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_NAME`: Database name
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `JWT_SECRET`: JWT signing secret
- `JWT_EXPIRY`: Token expiry duration
- `REDIS_URL`: Redis connection URL

## Data Models

### User Model
```go
type User struct {
    ID        int64     `json:"id"`
    Email     string    `json:"email"`
    Username  string    `json:"username"`
    Password  string    `json:"-"`
    FirstName string    `json:"firstName"`
    LastName  string    `json:"lastName"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

## Development

### Running API Service
```bash
cd apps/auth/api
go run main.go
```

### Running RPC Service
```bash
cd apps/auth/rpc
go run main.go
```

### Testing
```bash
go test ./tests/...
```

### Generating Protobuf
```bash
protoc --go_out=. --go-grpc_out=. pb/auth.proto
```

## Docker

### API Service
```bash
docker build -t smart-document/auth-api -f apps/auth/api/Dockerfile .
```

### RPC Service
```bash
docker build -t smart-document/auth-rpc -f apps/auth/rpc/Dockerfile .
```

## Security

- Passwords are hashed using bcrypt
- JWT tokens are signed with RS256
- Rate limiting on login attempts
- Input validation on all endpoints