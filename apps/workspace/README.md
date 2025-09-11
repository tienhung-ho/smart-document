# Workspace Service

Workspace and team management service for Smart Document platform.

## Overview

The Workspace service manages workspaces, team members, permissions, and templates. It provides both HTTP API and gRPC endpoints for workspace operations.

## Architecture

- **API Component**: HTTP endpoints for workspace management
- **RPC Component**: gRPC service for permission checks and workspace info

## Features

- Workspace creation and management
- Team member invitations and management
- Role-based access control (RBAC)
- Workspace templates
- Permission management
- Audit logging

## Components

### API Service (`/api`)
HTTP REST API for workspace operations.

**Endpoints:**
- `GET /workspaces` - List workspaces
- `POST /workspaces` - Create workspace
- `GET /workspaces/:id` - Get workspace
- `PUT /workspaces/:id` - Update workspace
- `DELETE /workspaces/:id` - Delete workspace
- `POST /workspaces/:id/invite` - Invite member
- `GET /workspaces/:id/members` - List members
- `PUT /workspaces/:id/members/:userId` - Update member role
- `DELETE /workspaces/:id/members/:userId` - Remove member

### RPC Service (`/rpc`)
gRPC service for internal communication.

**Methods:**
- `CheckWorkspacePermission(PermissionRequest) -> PermissionResponse`
- `GetWorkspaceInfo(WorkspaceRequest) -> WorkspaceResponse`
- `GetUserWorkspaces(UserRequest) -> WorkspacesResponse`

## Configuration

### Environment Variables

- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_NAME`: Database name
- `EMAIL_SERVICE_URL`: Email service for invitations
- `REDIS_URL`: Redis for caching

## Data Models

### Workspace Model
```go
type Workspace struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    OwnerID     int64     `json:"ownerId"`
    Plan        string    `json:"plan"`
    Settings    Settings  `json:"settings"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### Workspace Member Model
```go
type WorkspaceMember struct {
    ID          int64     `json:"id"`
    WorkspaceID int64     `json:"workspaceId"`
    UserID      int64     `json:"userId"`
    Role        string    `json:"role"` // owner, admin, editor, viewer
    Status      string    `json:"status"` // active, pending, suspended
    InvitedBy   int64     `json:"invitedBy"`
    JoinedAt    time.Time `json:"joinedAt"`
}
```

### Workspace Template Model
```go
type WorkspaceTemplate struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Category    string    `json:"category"`
    Config      Template  `json:"config"`
    IsPublic    bool      `json:"isPublic"`
    CreatedBy   int64     `json:"createdBy"`
}
```

## Roles and Permissions

### Roles
- **Owner**: Full access, can delete workspace
- **Admin**: Manage members, settings, documents
- **Editor**: Create and edit documents
- **Viewer**: Read-only access

### Permissions Matrix
| Action | Owner | Admin | Editor | Viewer |
|--------|-------|-------|--------|--------|
| Delete Workspace | ✓ | ✗ | ✗ | ✗ |
| Manage Members | ✓ | ✓ | ✗ | ✗ |
| Edit Settings | ✓ | ✓ | ✗ | ✗ |
| Create Documents | ✓ | ✓ | ✓ | ✗ |
| Edit Documents | ✓ | ✓ | ✓ | ✗ |
| View Documents | ✓ | ✓ | ✓ | ✓ |

## Templates

### Pre-built Templates
- **Startup**: Basic workspace for small teams
- **Enterprise**: Advanced features for large organizations
- **Education**: Templates for educational institutions
- **Non-profit**: Templates for non-profit organizations

## Development

### Running API Service
```bash
cd apps/workspace/api
go run main.go
```

### Running RPC Service
```bash
cd apps/workspace/rpc
go run main.go
```

### Testing
```bash
go test ./tests/...
```

## Docker

### API Service
```bash
docker build -t smart-document/workspace-api -f apps/workspace/api/Dockerfile .
```

### RPC Service
```bash
docker build -t smart-document/workspace-rpc -f apps/workspace/rpc/Dockerfile .
```

## Dependencies

- PostgreSQL for data storage
- Redis for caching
- Email service for invitations