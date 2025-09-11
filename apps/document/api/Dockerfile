# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy common dependencies
COPY common/ ./common/

# Copy document source
COPY apps/document/ ./apps/document/

# Build the application
WORKDIR /app/apps/document/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary
COPY --from=builder /app/apps/document/api/main .

# Copy config file
COPY --from=builder /app/apps/document/api/etc/ ./etc/

# Expose port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1

CMD ["./main"]