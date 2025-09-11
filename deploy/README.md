# Smart Document Deployment Guide

This guide covers deploying the Smart Document microservices platform using Docker and Kubernetes.

## Quick Start

### Using Docker Compose (Local Development)

1. **Start all services:**
```bash
docker-compose up -d
```

2. **View logs:**
```bash
docker-compose logs -f [service-name]
```

3. **Stop services:**
```bash
docker-compose down
```

### Using Kubernetes

1. **Create namespace and apply configs:**
```bash
kubectl apply -f deploy/k8s/namespace-config.yaml
```

2. **Deploy infrastructure services:**
```bash
kubectl apply -f deploy/k8s/database.yaml
kubectl apply -f deploy/k8s/storage.yaml
```

3. **Deploy application services:**
```bash
kubectl apply -f deploy/k8s/auth-service.yaml
kubectl apply -f deploy/k8s/document-service.yaml
kubectl apply -f deploy/k8s/collaboration-service.yaml
kubectl apply -f deploy/k8s/workspace-service.yaml
kubectl apply -f deploy/k8s/gateway.yaml
```

4. **Setup ingress:**
```bash
kubectl apply -f deploy/k8s/ingress.yaml
```

## Service Endpoints

### Local Development (Docker Compose)
- Gateway: http://localhost:8080
- Auth API: http://localhost:8081
- Document API: http://localhost:8082
- Collaboration API: http://localhost:8083
- Workspace API: http://localhost:8084
- MinIO Console: http://localhost:9001

### Kubernetes (with Ingress)
- Main App: http://smart-document.local
- API Gateway: http://api.smart-document.local
- MinIO Console: http://minio.smart-document.local

## Environment Variables

### Required for all services:
- `DB_HOST`: Database host
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `REDIS_HOST`: Redis host
- `JWT_SECRET`: JWT signing secret

### Additional for Document Service:
- `ELASTICSEARCH_URL`: Elasticsearch URL
- `STORAGE_ENDPOINT`: MinIO/S3 endpoint
- `STORAGE_ACCESS_KEY`: Storage access key
- `STORAGE_SECRET_KEY`: Storage secret key

### Additional for Collaboration Service:
- `KAFKA_BROKERS`: Kafka broker URLs

## Scaling

### Docker Compose
```bash
docker-compose up -d --scale auth-api=3 --scale document-api=2
```

### Kubernetes
```bash
kubectl scale deployment auth-api --replicas=3 -n smart-document
kubectl scale deployment document-api --replicas=2 -n smart-document
```

## Monitoring

### Health Checks
- Gateway: `GET /health`
- Each API service: `GET /health`
- gRPC services: Use grpc_health_probe

### Logs
```bash
# Docker Compose
docker-compose logs -f [service-name]

# Kubernetes
kubectl logs -f deployment/[service-name] -n smart-document
```

## Troubleshooting

### Common Issues

1. **Database connection failed:**
   - Check PostgreSQL is running
   - Verify connection credentials
   - Ensure database exists

2. **Redis connection failed:**
   - Check Redis is running
   - Verify Redis host/port

3. **File upload fails:**
   - Check MinIO is running
   - Verify storage credentials
   - Ensure bucket exists

4. **Real-time collaboration not working:**
   - Check Kafka is running
   - Verify WebSocket connections
   - Check Redis for session storage

### Debug Commands

```bash
# Check service status
kubectl get pods -n smart-document

# View service logs
kubectl logs -f pod/[pod-name] -n smart-document

# Port forward for debugging
kubectl port-forward svc/postgres-service 5432:5432 -n smart-document
```