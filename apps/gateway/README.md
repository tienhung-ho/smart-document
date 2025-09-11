# Gateway Service

API Gateway service for Smart Document microservices.

## Overview

The Gateway service acts as the single entry point for all client requests, routing them to appropriate microservices and handling cross-cutting concerns like authentication, rate limiting, and CORS.

## Features

- Request routing to microservices
- Authentication middleware
- Rate limiting
- CORS handling
- WebSocket proxy support
- Request/response logging

## Configuration

Configuration file: `etc/gateway.yaml`

### Environment Variables

- `PORT`: Server port (default: 8080)
- `AUTH_SERVICE_URL`: Authentication service URL
- `DOCUMENT_SERVICE_URL`: Document service URL
- `COLLABORATION_SERVICE_URL`: Collaboration service URL
- `WORKSPACE_SERVICE_URL`: Workspace service URL

## API Routes

### Authentication Routes
- `POST /api/auth/login`
- `POST /api/auth/register`
- `POST /api/auth/refresh`
- `GET /api/auth/profile`

### Document Routes
- `GET /api/documents`
- `POST /api/documents`
- `GET /api/documents/:id`
- `PUT /api/documents/:id`
- `DELETE /api/documents/:id`

### Collaboration Routes
- `WS /api/collaboration/ws/:documentId`
- `GET /api/collaboration/presence/:documentId`

### Workspace Routes
- `GET /api/workspaces`
- `POST /api/workspaces`
- `GET /api/workspaces/:id`

## Development

### Running Locally

```bash
cd apps/gateway
go run main.go
```

### Testing

```bash
go test ./tests/...
```

### Building

```bash
go build -o gateway main.go
```

## Docker

Build image:
```bash
docker build -t smart-document/gateway .
```

Run container:
```bash
docker run -p 8080:8080 smart-document/gateway
```

## Health Check

- `GET /health` - Service health status