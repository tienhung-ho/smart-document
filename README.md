# Smart Document Microservice

A robust document management microservice built with Go, featuring hybrid configuration, centralized error handling, and structured logging.

## 🏗️ Architecture

This project follows a microservice architecture with the following services:

- **Gateway**: API Gateway with routing and middleware
- **Auth**: Authentication and authorization service (API + RPC)
- **Document**: Document CRUD and file management (API + RPC)
- **Collaboration**: Real-time editing and presence (API + RPC)
- **Workspace**: Team management and permissions (API + RPC)

## 🛠️ Core Features

### Hybrid Configuration (Viper)
- YAML configuration files with environment-specific overrides
- Environment variable support with prefix `SD_`
- Hierarchical configuration loading
- Default values for all settings

### Error Handling
- Centralized error management with custom error codes
- Structured error context and details
- HTTP status code mapping
- Stack trace collection for debugging

### Logging (Zap)
- High-performance structured logging
- Multiple output formats (JSON, Console)
- File rotation with lumberjack
- Context-aware logging
- HTTP request and database query logging

## 📁 Project Structure

```
smart-document/
├── apps/                    # Microservices
│   ├── gateway/            # API Gateway
│   ├── auth/               # Authentication Service
│   ├── document/           # Document Management
│   ├── collaboration/      # Real-time Collaboration
│   └── workspace/          # Workspace Management
├── common/                 # Shared Libraries
│   ├── config/            # Configuration management
│   ├── errors/            # Error handling
│   ├── logging/           # Logging utilities
│   ├── middleware/        # HTTP middleware
│   ├── utils/             # Common utilities
│   └── database/          # Database utilities
└── migrations/            # Database migrations
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL
- Redis
- Kafka (optional)

### Setup

1. **Clone and setup**:
   ```bash
   git clone <repository>
   cd smart-document
   go mod download
   ```

2. **Configure environment**:
   ```bash
   cp apps/gateway/etc/gateway.yaml apps/gateway/etc/config-local.yaml
   # Edit config-local.yaml with your settings
   ```

3. **Run the gateway service**:
   ```bash
   cd apps/gateway
   go run main.go
   ```

## ⚙️ Configuration

### Configuration Files
- `gateway.yaml` - Default configuration
- `config-production.yaml` - Production overrides
- `config-development.yaml` - Development overrides

### Environment Variables
All configuration can be overridden with environment variables using the `SD_` prefix:

```bash
export SD_ENVIRONMENT=production
export SD_SERVER_PORT=8080
export SD_DATABASE_HOST=localhost
export SD_JWT_SECRET=your-secret-key
```

### Configuration Structure
```yaml
environment: development
server:
  host: 0.0.0.0
  port: 8080
database:
  driver: postgres
  host: localhost
  port: 5432
logging:
  level: info
  format: json
  output: both
```

## 🔍 Error Handling

The project includes a comprehensive error handling system:

```go
// Create custom errors
err := errors.Validation("Invalid email format").
    WithContext("field", "email").
    WithDetails("Email must be in valid format")

// Handle errors
if appErr, ok := errors.GetAppError(err); ok {
    logging.Errorf("Error [%d]: %s", appErr.Code, appErr.Message)
    // Returns appropriate HTTP status code
    return c.Status(appErr.HTTPStatus()).JSON(appErr)
}
```

## 📝 Logging

Structured logging with multiple output options:

```go
// Simple logging
logging.Info("Service started")
logging.Errorf("Failed to connect: %v", err)

// Structured logging
logging.WithFields(map[string]interface{}{
    "user_id": "123",
    "action": "login",
}).Info("User action")

// HTTP request logging
logging.LogHTTPRequest("POST", "/api/login", "user123", 200, 150.5)
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific service tests
go test ./apps/auth/...
```

## 🐳 Docker

```bash
# Build all services
make docker-build

# Start all services
make docker-up

# Stop all services
make docker-down
```

## 📊 Monitoring

The project includes monitoring setup:
- Prometheus metrics
- Grafana dashboards
- Jaeger tracing
- Health check endpoints

## 🔧 Development

### Adding a New Service

1. Create service directory structure:
   ```bash
   mkdir -p apps/myservice/{api,rpc,model,tests}
   ```

2. Create configuration:
   ```bash
   cp apps/gateway/etc/gateway.yaml apps/myservice/etc/myservice.yaml
   ```

3. Update configuration values for your service

4. Create main.go following the gateway example

### Common Patterns

**Loading Configuration:**
```go
cfg, err := config.LoadConfig("./etc", "service-name")
```

**Initializing Logger:**
```go
logger, err := logging.InitLogger(&cfg.Logging)
```

**Error Handling:**
```go
return errors.BadRequest("Invalid input").WithContext("field", value)
```

## 📜 License

This project is licensed under the MIT License.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📞 Support

For support and questions, please open an issue in the repository.