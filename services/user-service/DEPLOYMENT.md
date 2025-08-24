# User Service - Deployment Guide

## 🚀 Production Deployment

### Prerequisites

1. **PostgreSQL Database**
   - Version 15+
   - Database `lunary_db` created
   - User with full access to the database

2. **Auth Service**
   - Running and accessible
   - Same JWT_SECRET configured

3. **Environment Configuration**
   - All required environment variables set
   - Network connectivity between services

### Environment Variables

#### Required
```bash
# Service
PORT=8083

# Database (REQUIRED)
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=lunary_db

# JWT (MUST match auth-service)
JWT_SECRET=your-production-jwt-secret
```

#### Optional
```bash
# Database
DB_SSL_MODE=require  # for production

# Integration
AUTH_SERVICE_URL=http://auth-service:8081

# Logging
LOG_LEVEL=info
```

### Database Setup

1. **Apply Migrations**
```bash
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f migrations/001_create_user_profiles_table.up.sql
```

2. **Verify Tables**
```sql
-- Check if tables exist
\dt

-- Verify foreign key constraint
\d user_profiles
```

### Docker Deployment

#### Build Image
```bash
docker build -t user-service:latest .
```

#### Run Container
```bash
docker run -d \
  --name user-service \
  -p 8083:8083 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-db-password \
  -e JWT_SECRET=your-jwt-secret \
  --restart unless-stopped \
  user-service:latest
```

#### Health Check
```bash
curl http://localhost:8083/health
```

### Kubernetes Deployment

#### ConfigMap
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: user-service-config
data:
  PORT: "8083"
  DB_HOST: "postgres-service"
  DB_PORT: "5432"
  DB_USER: "postgres"
  DB_NAME: "lunary_db"
  DB_SSL_MODE: "require"
  AUTH_SERVICE_URL: "http://auth-service:8081"
```

#### Secret
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: user-service-secrets
type: Opaque
stringData:
  DB_PASSWORD: "your-db-password"
  JWT_SECRET: "your-jwt-secret"
```

#### Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  labels:
    app: user-service
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: user-service:latest
        ports:
        - containerPort: 8083
          name: http
        envFrom:
        - configMapRef:
            name: user-service-config
        - secretRef:
            name: user-service-secrets
        
        # Health Checks
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8083
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8083
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        
        # Resource Limits
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        
        # Security
        securityContext:
          allowPrivilegeEscalation: false
          runAsNonRoot: true
          runAsUser: 1000
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
```

#### Service
```yaml
apiVersion: v1
kind: Service
metadata:
  name: user-service
  labels:
    app: user-service
spec:
  selector:
    app: user-service
  ports:
  - port: 8083
    targetPort: 8083
    name: http
  type: ClusterIP
```

#### Ingress (Optional)
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: user-service-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: api.lunary.com
    http:
      paths:
      - path: /user
        pathType: Prefix
        backend:
          service:
            name: user-service
            port:
              number: 8083
```

### Load Balancer Configuration

#### NGINX
```nginx
upstream user_service {
    server user-service-1:8083;
    server user-service-2:8083;
    server user-service-3:8083;
}

server {
    listen 80;
    server_name api.lunary.com;

    location /api/v1/profile {
        proxy_pass http://user_service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Health check
        proxy_connect_timeout 5s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }

    location /health {
        proxy_pass http://user_service;
        access_log off;
    }
}
```

### Monitoring Setup

#### Prometheus Metrics (Future)
```yaml
apiVersion: v1
kind: Service
metadata:
  name: user-service-metrics
  labels:
    app: user-service
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "8083"
    prometheus.io/path: "/metrics"
spec:
  selector:
    app: user-service
  ports:
  - port: 8083
    name: metrics
```

#### Grafana Dashboard Query Examples
```promql
# Request rate
rate(http_requests_total{service="user-service"}[5m])

# Error rate
rate(http_requests_total{service="user-service",status=~"5.."}[5m])

# Response time
histogram_quantile(0.95, http_request_duration_seconds{service="user-service"})

# Database connections
db_connections_active{service="user-service"}
```

### Backup Strategy

#### Database Backup
```bash
# Daily backup
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME \
  --no-password --clean --if-exists \
  > user_service_backup_$(date +%Y%m%d).sql

# Restore
psql -h $DB_HOST -U $DB_USER -d $DB_NAME < backup_file.sql
```

### Security Considerations

1. **Database Security**
   - Use SSL connections in production (DB_SSL_MODE=require)
   - Rotate database passwords regularly
   - Limit database user permissions

2. **JWT Security**
   - Use strong, unique JWT secrets
   - Sync JWT secrets with auth-service
   - Consider JWT rotation strategy

3. **Network Security**
   - Use internal networks for service communication
   - Implement proper firewall rules
   - Consider service mesh for mTLS

4. **Container Security**
   - Run as non-root user
   - Use read-only root filesystem
   - Drop all capabilities
   - Scan images for vulnerabilities

### Troubleshooting

#### Common Issues

1. **Database Connection Failed**
```bash
# Check connectivity
telnet $DB_HOST $DB_PORT

# Check credentials
psql -h $DB_HOST -U $DB_USER -d $DB_NAME

# Check logs
kubectl logs deployment/user-service
```

2. **JWT Token Invalid**
```bash
# Verify JWT secret matches auth-service
echo $JWT_SECRET

# Test token from auth-service
curl -H "Authorization: Bearer $TOKEN" http://localhost:8083/api/v1/profile
```

3. **Health Check Failed**
```bash
# Check health endpoints
curl http://localhost:8083/health
curl http://localhost:8083/health/live
curl http://localhost:8083/health/ready
```

#### Log Analysis
```bash
# Check application logs
kubectl logs -f deployment/user-service

# Check specific errors
kubectl logs deployment/user-service | grep ERROR

# Check database operations
kubectl logs deployment/user-service | grep "database"
```

### Performance Tuning

#### Database Optimization
```sql
-- Monitor slow queries
SELECT query, mean_time, calls 
FROM pg_stat_statements 
WHERE query LIKE '%user_profiles%'
ORDER BY mean_time DESC;

-- Check index usage
SELECT schemaname, tablename, attname, n_distinct, correlation
FROM pg_stats
WHERE tablename = 'user_profiles';
```

#### Application Tuning
```bash
# Monitor resource usage
kubectl top pods -l app=user-service

# Check memory/CPU limits
kubectl describe pods -l app=user-service
```

### Scaling Guidelines

#### Horizontal Scaling
- 1 replica: < 100 RPS
- 3 replicas: 100-500 RPS  
- 5+ replicas: > 500 RPS

#### Database Scaling
- Connection pool tuning based on replica count
- Read replicas for read-heavy workloads
- Connection pooling (PgBouncer) for high connection counts

### Rollback Procedure

1. **Application Rollback**
```bash
# Kubernetes
kubectl rollout undo deployment/user-service

# Docker
docker stop user-service
docker run -d --name user-service user-service:previous-tag
```

2. **Database Rollback**
```bash
# Apply down migration
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f migrations/001_create_user_profiles_table.down.sql
```

### Maintenance

#### Regular Tasks
- [ ] Monitor health check status
- [ ] Review application logs for errors
- [ ] Check database performance
- [ ] Update dependencies (monthly)
- [ ] Rotate secrets (quarterly)
- [ ] Test backup/restore procedure (quarterly)

#### Emergency Procedures
- [ ] Have rollback plan ready
- [ ] Monitor key metrics during deployment
- [ ] Keep database backups current
- [ ] Have emergency contact information available
