# Document Service

Document management service for Smart Document platform.

## Overview

The Document service handles document CRUD operations, file uploads, version management, and search functionality. It provides both HTTP API and gRPC endpoints.

## Architecture

- **API Component**: HTTP endpoints for document operations
- **RPC Component**: gRPC service for internal service communication

## Features

- Document CRUD operations
- File upload and storage
- Document versioning and history
- Full-text search with Elasticsearch
- Document permissions and sharing
- Metadata management
- File format conversion

## Components

### API Service (`/api`)
HTTP REST API for document operations.

**Endpoints:**
- `GET /documents` - List documents
- `POST /documents` - Create document
- `GET /documents/:id` - Get document
- `PUT /documents/:id` - Update document
- `DELETE /documents/:id` - Delete document
- `POST /documents/:id/upload` - Upload file
- `GET /documents/:id/versions` - Version history
- `POST /documents/:id/share` - Share document

### RPC Service (`/rpc`)
gRPC service for internal communication.

**Methods:**
- `GetDocument(DocumentRequest) -> DocumentResponse`
- `CheckPermission(PermissionRequest) -> PermissionResponse`
- `SearchDocuments(SearchRequest) -> SearchResponse`

## Configuration

### Environment Variables

- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_NAME`: Database name
- `STORAGE_TYPE`: Storage backend (minio/s3)
- `STORAGE_ENDPOINT`: Storage endpoint
- `STORAGE_ACCESS_KEY`: Storage access key
- `STORAGE_SECRET_KEY`: Storage secret key
- `ELASTICSEARCH_URL`: Elasticsearch cluster URL

## Data Models

### Document Model
```go
type Document struct {
    ID          int64     `json:"id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    ContentType string    `json:"contentType"`
    Size        int64     `json:"size"`
    AuthorID    int64     `json:"authorId"`
    WorkspaceID int64     `json:"workspaceId"`
    Status      string    `json:"status"`
    Tags        []string  `json:"tags"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### Document Version Model
```go
type DocumentVersion struct {
    ID         int64     `json:"id"`
    DocumentID int64     `json:"documentId"`
    Version    int       `json:"version"`
    Content    string    `json:"content"`
    Changes    string    `json:"changes"`
    AuthorID   int64     `json:"authorId"`
    CreatedAt  time.Time `json:"createdAt"`
}
```

## Development

### Running API Service
```bash
cd apps/document/api
go run main.go
```

### Running RPC Service
```bash
cd apps/document/rpc
go run main.go
```

### Testing
```bash
go test ./tests/...
```

## Docker

### API Service
```bash
docker build -t smart-document/document-api -f apps/document/api/Dockerfile .
```

### RPC Service
```bash
docker build -t smart-document/document-rpc -f apps/document/rpc/Dockerfile .
```

## Dependencies

- PostgreSQL for metadata storage
- MinIO/S3 for file storage
- Elasticsearch for search functionality
- Redis for caching