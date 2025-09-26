# Deployment Guide

This guide covers all deployment options for the SIL Backend Assessment application.

## Table of Contents

1. [Local Development](#local-development)
2. [Docker Deployment](#docker-deployment)
3. [Kubernetes Deployment (Minikube)](#kubernetes-deployment-minikube)
4. [Production Deployment](#production-deployment)
5. [Monitoring and Health Checks](#monitoring-and-health-checks)
6. [Troubleshooting](#troubleshooting)

## Local Development

### Prerequisites
- Go 1.24+
- PostgreSQL 12+
- Redis 6+ (optional)

### Setup
```bash
# Clone repository
git clone <repository-url>
cd sil-backend-assessment

# Install dependencies
make deps

# Setup database
createdb sil_backend_assessment_db
make migrate_up

# Run application
make build
./bin/server
```

### Access Points
- **API**: http://localhost:8080/api
- **GraphQL**: http://localhost:8080/graphql
- **Health**: http://localhost:8080/api/health
- **Docs**: http://localhost:8080/docs

## Docker Deployment

### Single Container
```bash
# Build image
make docker-build

# Run container
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PASSWORD=your_password \
  sil-backend-assessment:latest
```

### Docker Compose
```bash
# Start all services
make docker-up

# Stop services
make docker-down
```

### Services Included
- **Application**: Main API server
- **PostgreSQL**: Database
- **Redis**: Caching layer
- **Prometheus**: Metrics collection

## Kubernetes Deployment (Minikube)

### Prerequisites
- Minikube installed
- kubectl configured
- Docker Desktop (for building images)

### Quick Deployment
```bash
# Deploy everything
make deploy-minikube

# Get access information
minikube service sil-backend-nodeport -n sil-backend-assessment --url
```

### Manual Deployment
```bash
# Start Minikube
minikube start --memory=4096 --cpus=2

# Set Docker environment
eval $(minikube docker-env)

# Build and load image
make docker-build
minikube image load sil-backend-assessment:latest

# Deploy to Kubernetes
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/redis.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
kubectl apply -f k8s/hpa.yaml
```

### Access Points
```bash
# Get Minikube IP
MINIKUBE_IP=$(minikube ip)

# Access URLs
echo "Health Check: http://${MINIKUBE_IP}:30080/api/health"
echo "API Docs: http://${MINIKUBE_IP}:30080/docs"
echo "GraphQL Playground: http://${MINIKUBE_IP}:30080/graphql/playground"
```

### Scaling
```bash
# Manual scaling
kubectl scale deployment sil-backend-assessment --replicas=5 -n sil-backend-assessment

# Check HPA status
kubectl get hpa -n sil-backend-assessment
```

### Cleanup
```bash
# Delete everything
make delete-minikube

# Or manual cleanup
kubectl delete namespace sil-backend-assessment
```

## Production Deployment

### Environment Configuration

#### ConfigMap Updates
```bash
kubectl patch configmap sil-backend-config -n sil-backend-assessment --patch '
data:
  config.yaml: |
    server:
      rest_port: 8080
      shutdown_timeout: 30s
    database:
      host: postgres-service
      port: 5432
      user: sil_user
      password: sil_password
      dbname: sil_backend_assessment_db
      sslmode: require
    redis:
      address: redis-service:6379
      password: "secure-redis-password"
      db: 0
    auth:
      jwt_secret: "your-production-jwt-secret"
      jwt_expiry: 1h
      jwt_refresh_secret: "your-production-refresh-secret"
      jwt_refresh_expiry: 7d
    smtp:
      host: smtp.your-domain.com
      port: 587
      username: noreply@your-domain.com
      password: "your-smtp-password"
      from: noreply@your-domain.com
      tls: true
    at:
      api_key: "your-production-at-api-key"
      username: "your-production-at-username"
      base_url: "https://api.africastalking.com"
    oidc:
      enabled: true
      provider_url: "https://your-oidc-provider.com"
      client_id: "your-production-client-id"
      client_secret: "your-production-client-secret"
      redirect_url: "https://your-domain.com/auth/oidc/callback"
      scopes:
        - openid
        - profile
        - email
    logging:
      level: info
      file: /var/log/sil-backend-assessment.log
    metrics:
      enabled: true
      port: 9090
      path: /metrics
'
```

#### Secrets Management
```bash
# Create production secrets
kubectl create secret generic sil-backend-secrets \
  --from-literal=jwt-secret="your-production-jwt-secret" \
  --from-literal=jwt-refresh-secret="your-production-refresh-secret" \
  --from-literal=smtp-password="your-production-smtp-password" \
  --from-literal=at-api-key="your-production-at-api-key" \
  --from-literal=oidc-client-secret="your-production-oidc-secret" \
  --from-literal=db-password="your-production-db-password" \
  --namespace=sil-backend-assessment
```

### Production Considerations

#### Security
- Use TLS/SSL for all communications
- Enable RBAC in Kubernetes
- Use Network Policies
- Regular security updates
- Secrets rotation

#### Performance
- Configure resource limits
- Enable HPA (Horizontal Pod Autoscaler)
- Use CDN for static assets
- Database connection pooling
- Redis caching strategy

#### Monitoring
- Prometheus metrics collection
- Grafana dashboards
- Log aggregation (ELK stack)
- Alerting (AlertManager)
- Health checks and probes

#### Backup
- Database backups
- Configuration backups
- Disaster recovery plan
- Point-in-time recovery

## Monitoring and Health Checks

### Health Check Endpoints
```bash
# Basic health check
curl http://localhost:8080/api/health

# Detailed health check (if implemented)
curl http://localhost:8080/api/health/detailed

# Readiness probe
curl http://localhost:8080/api/ready

# Liveness probe
curl http://localhost:8080/api/live
```

### Metrics Endpoint
```bash
# Prometheus metrics
curl http://localhost:9090/metrics
```

### Kubernetes Probes
```yaml
livenessProbe:
  httpGet:
    path: /api/health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /api/health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
```

### Log Monitoring
```bash
# View application logs
kubectl logs -f deployment/sil-backend-assessment -n sil-backend-assessment

# View logs from all pods
kubectl logs -f -l app=sil-backend-assessment -n sil-backend-assessment

# View logs with timestamps
kubectl logs -f deployment/sil-backend-assessment -n sil-backend-assessment --timestamps
```

## Troubleshooting

### Common Issues

#### Application Won't Start
```bash
# Check pod status
kubectl get pods -n sil-backend-assessment

# Check pod logs
kubectl logs <pod-name> -n sil-backend-assessment

# Check pod events
kubectl describe pod <pod-name> -n sil-backend-assessment
```

#### Database Connection Issues
```bash
# Check database pod
kubectl get pods -l app=postgres -n sil-backend-assessment

# Check database logs
kubectl logs -l app=postgres -n sil-backend-assessment

# Test database connection
kubectl exec -it <postgres-pod> -n sil-backend-assessment -- psql -U sil_user -d sil_backend_assessment_db
```

#### Service Discovery Issues
```bash
# Check services
kubectl get services -n sil-backend-assessment

# Check endpoints
kubectl get endpoints -n sil-backend-assessment

# Test service connectivity
kubectl exec -it <app-pod> -n sil-backend-assessment -- curl postgres-service:5432
```

#### Resource Issues
```bash
# Check resource usage
kubectl top pods -n sil-backend-assessment

# Check node resources
kubectl top nodes

# Check HPA status
kubectl get hpa -n sil-backend-assessment
```

### Debug Commands

#### Application Debugging
```bash
# Port forward to access application locally
kubectl port-forward svc/sil-backend-service 8080:8080 -n sil-backend-assessment

# Access application
curl http://localhost:8080/api/health
```

#### Database Debugging
```bash
# Port forward to database
kubectl port-forward svc/postgres-service 5432:5432 -n sil-backend-assessment

# Connect to database
psql -h localhost -p 5432 -U sil_user -d sil_backend_assessment_db
```

#### Network Debugging
```bash
# Test internal DNS
kubectl exec -it <app-pod> -n sil-backend-assessment -- nslookup postgres-service

# Test connectivity
kubectl exec -it <app-pod> -n sil-backend-assessment -- ping postgres-service
```

### Performance Optimization

#### Resource Tuning
```yaml
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "1Gi"
    cpu: "1000m"
```

#### Scaling Configuration
```yaml
spec:
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Backup and Recovery

#### Database Backup
```bash
# Create backup
kubectl exec <postgres-pod> -n sil-backend-assessment -- pg_dump -U sil_user sil_backend_assessment_db > backup.sql

# Restore backup
kubectl exec -i <postgres-pod> -n sil-backend-assessment -- psql -U sil_user sil_backend_assessment_db < backup.sql
```

#### Configuration Backup
```bash
# Backup ConfigMaps
kubectl get configmap sil-backend-config -n sil-backend-assessment -o yaml > config-backup.yaml

# Backup Secrets
kubectl get secret sil-backend-secrets -n sil-backend-assessment -o yaml > secrets-backup.yaml
```

## Support

For deployment issues:
1. Check the troubleshooting section above
2. Review application logs
3. Verify configuration
4. Test connectivity between services
5. Check resource utilization

For additional help, please refer to the main README.md or create an issue in the repository.
