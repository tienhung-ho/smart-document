# Collaboration Service

Real-time collaboration service for Smart Document platform.

## Overview

The Collaboration service enables real-time collaborative editing using Operational Transformation (OT), presence tracking, and cursor synchronization. It provides WebSocket connections and gRPC endpoints.

## Architecture

- **API Component**: WebSocket handlers for real-time communication
- **RPC Component**: gRPC service for collaboration session management

## Features

- Real-time collaborative editing
- Operational Transformation (OT) engine
- User presence tracking
- Cursor position synchronization
- Conflict resolution
- Session management
- Event broadcasting

## Components

### API Service (`/api`)
WebSocket API for real-time collaboration.

**Endpoints:**
- `WS /ws/:documentId` - WebSocket connection for document
- `GET /presence/:documentId` - Get active users
- `POST /sessions` - Create collaboration session

### RPC Service (`/rpc`)
gRPC service for session management.

**Methods:**
- `ApplyOperation(OperationRequest) -> OperationResponse`
- `GetSession(SessionRequest) -> SessionResponse`
- `JoinSession(JoinRequest) -> JoinResponse`
- `LeaveSession(LeaveRequest) -> LeaveResponse`

## Configuration

### Environment Variables

- `REDIS_URL`: Redis connection for session storage
- `KAFKA_BROKERS`: Kafka brokers for event streaming
- `WS_PORT`: WebSocket server port
- `RPC_PORT`: gRPC server port

## Data Models

### Operation Model
```go
type Operation struct {
    ID         string    `json:"id"`
    Type       string    `json:"type"` // insert, delete, retain
    Position   int       `json:"position"`
    Content    string    `json:"content,omitempty"`
    Length     int       `json:"length,omitempty"`
    AuthorID   int64     `json:"authorId"`
    Timestamp  time.Time `json:"timestamp"`
    DocumentID int64     `json:"documentId"`
}
```

### Session Model
```go
type Session struct {
    ID          string    `json:"id"`
    DocumentID  int64     `json:"documentId"`
    Users       []User    `json:"users"`
    Operations  []Operation `json:"operations"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### Presence Model
```go
type Presence struct {
    UserID     int64     `json:"userId"`
    DocumentID int64     `json:"documentId"`
    CursorPos  int       `json:"cursorPosition"`
    Selection  Selection `json:"selection"`
    LastSeen   time.Time `json:"lastSeen"`
}
```

## WebSocket Events

### Client to Server
- `join_document` - Join collaboration session
- `operation` - Send operation
- `cursor_move` - Update cursor position
- `selection_change` - Update text selection

### Server to Client
- `operation_applied` - Operation was applied
- `user_joined` - User joined session
- `user_left` - User left session
- `cursor_update` - Cursor position update
- `presence_update` - User presence update

## Operational Transformation

The service implements OT for conflict resolution:

1. **Transform Operations**: Adjust operations based on concurrent changes
2. **Apply Operations**: Apply transformed operations to document
3. **Broadcast Changes**: Send updates to all connected clients
4. **Maintain Consistency**: Ensure all clients have same document state

## Development

### Running API Service
```bash
cd apps/collaboration/api
go run main.go
```

### Running RPC Service
```bash
cd apps/collaboration/rpc
go run main.go
```

### Testing
```bash
go test ./tests/...
```

## Docker

### API Service
```bash
docker build -t smart-document/collaboration-api -f apps/collaboration/api/Dockerfile .
```

### RPC Service
```bash
docker build -t smart-document/collaboration-rpc -f apps/collaboration/rpc/Dockerfile .
```

## Dependencies

- Redis for session storage and pub/sub
- Kafka for event streaming
- WebSocket for real-time communication