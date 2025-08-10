# Production Deployment Guide

## 🎯 Обзор

Это руководство описывает deployment LUNARY в production environment с акцентом на reliability, security и observability.

## ✅ Pre-Production Checklist

### Code Quality
- [x] **Unit Tests** - 100% покрытие business logic
- [x] **Integration Tests** - database, HTTP endpoints
- [x] **Error Handling** - structured errors с HTTP codes
- [x] **Input Validation** - go-playground/validator
- [x] **Security Review** - SQL injection, XSS protection
- [x] **Performance Testing** - load testing, memory profiling

### Infrastructure
- [x] **Health Checks** - liveness, readiness, startup probes
- [x] **Graceful Shutdown** - 30s timeout с cleanup
- [x] **Logging** - structured JSON logs
- [x] **Monitoring** - metrics collection готовность
- [x] **Database** - migrations, backups, indexing
- [x] **Secrets Management** - Kubernetes secrets

## 🚀 Production Deployment

### 1. Environment Setup

#### Database Setup
```sql
-- Create production database
CREATE DATABASE lunary_production;
CREATE USER lunary_prod WITH ENCRYPTED PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE lunary_production TO lunary_prod;

-- Enable extensions
\c lunary_production
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";
```

#### Kubernetes Namespace
```bash
# Create namespace
kubectl create namespace lunary-prod

# Create secrets
kubectl create secret generic lunary-secrets \
  --from-literal=db-user=lunary_prod \
  --from-literal=db-password=secure_password \
  --from-literal=jwt-secret=super_secure_jwt_secret_change_me \
  -n lunary-prod
```

### 2. Database Deployment

```yaml
# postgres-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: lunary-prod
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    spec:
      containers:
      - name: postgres
        image: postgres:15
        env:
        - name: POSTGRES_DB
          value: "lunary_production"
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: lunary-secrets
              key: db-user
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: lunary-secrets
              key: db-password
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        livenessProbe:
          exec:
            command:
            - pg_isready
            - -U
            - lunary_prod
            - -d
            - lunary_production
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          exec:
            command:
            - pg_isready
            - -U
            - lunary_prod
            - -d
            - lunary_production
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc
```

### 3. Auth Service Deployment

```yaml
# auth-service-prod.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: lunary-prod
  labels:
    app: auth-service
    version: v1.0.0
spec:
  replicas: 3  # High availability
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
        version: v1.0.0
    spec:
      terminationGracePeriodSeconds: 45
      containers:
      - name: auth-service
        image: lunary/auth-service:v1.0.0
        ports:
        - containerPort: 8081
          name: http
        env:
        - name: DB_HOST
          value: "postgres-service"
        - name: DB_PORT
          value: "5432"
        - name: DB_NAME
          value: "lunary_production"
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: lunary-secrets
              key: db-user
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: lunary-secrets
              key: db-password
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: lunary-secrets
              key: jwt-secret
        - name: LOG_LEVEL
          value: "info"
        - name: PORT
          value: "8081"
        
        # Health checks
        livenessProbe:
          httpGet:
            path: /health/live
            port: http
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        
        readinessProbe:
          httpGet:
            path: /health/ready
            port: http
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        
        startupProbe:
          httpGet:
            path: /health/live
            port: http
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 30
        
        # Graceful shutdown
        lifecycle:
          preStop:
            exec:
              command: ["/bin/sh", "-c", "sleep 5"]
        
        # Resource limits
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
        
        # Security context
        securityContext:
          runAsNonRoot: true
          runAsUser: 1000
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL

---
apiVersion: v1
kind: Service
metadata:
  name: auth-service
  namespace: lunary-prod
spec:
  selector:
    app: auth-service
  ports:
  - port: 8081
    targetPort: 8081
    name: http
  type: ClusterIP
```

### 4. Load Balancer & Ingress

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: lunary-ingress
  namespace: lunary-prod
  annotations:
    kubernetes.io/ingress.class: "nginx"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  tls:
  - hosts:
    - api.lunary.app
    secretName: lunary-tls
  rules:
  - host: api.lunary.app
    http:
      paths:
      - path: /api/v1/auth
        pathType: Prefix
        backend:
          service:
            name: auth-service
            port:
              number: 8081
      - path: /health
        pathType: Prefix
        backend:
          service:
            name: auth-service
            port:
              number: 8081
```

### 5. Horizontal Pod Autoscaler

```yaml
# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auth-service-hpa
  namespace: lunary-prod
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auth-service
  minReplicas: 3
  maxReplicas: 10
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
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
```

## 📊 Monitoring Setup

### 1. Service Monitor (Prometheus)

```yaml
# service-monitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: auth-service-monitor
  namespace: lunary-prod
spec:
  selector:
    matchLabels:
      app: auth-service
  endpoints:
  - port: http
    path: /metrics
    interval: 30s
    scrapeTimeout: 10s
```

### 2. Grafana Dashboard

```json
{
  "dashboard": {
    "title": "LUNARY Auth Service",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total{service=\"auth-service\"}[5m])"
          }
        ]
      },
      {
        "title": "Error Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total{service=\"auth-service\",status=~\"5..\"}[5m]) / rate(http_requests_total{service=\"auth-service\"}[5m])"
          }
        ]
      },
      {
        "title": "Response Time",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket{service=\"auth-service\"}[5m]))"
          }
        ]
      }
    ]
  }
}
```

### 3. Alerting Rules

```yaml
# alerts.yaml
groups:
- name: lunary-auth-service
  rules:
  - alert: AuthServiceDown
    expr: up{job="auth-service"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Auth Service is down"
      description: "Auth service has been down for more than 1 minute"
  
  - alert: HighErrorRate
    expr: rate(http_requests_total{service="auth-service",status=~"5.."}[5m]) / rate(http_requests_total{service="auth-service"}[5m]) > 0.05
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High error rate on Auth Service"
      description: "Error rate is {{ $value | humanizePercentage }}"
  
  - alert: HighResponseTime
    expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket{service="auth-service"}[5m])) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High response time on Auth Service"
      description: "95th percentile response time is {{ $value }}s"
```

## 🔒 Security Configuration

### 1. Network Policies

```yaml
# network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: auth-service-netpol
  namespace: lunary-prod
spec:
  podSelector:
    matchLabels:
      app: auth-service
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    - podSelector:
        matchLabels:
          app: api-gateway
    ports:
    - protocol: TCP
      port: 8081
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - protocol: TCP
      port: 5432
  - to: []  # Allow DNS
    ports:
    - protocol: UDP
      port: 53
```

### 2. Pod Security Policy

```yaml
# pod-security-policy.yaml
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: lunary-psp
spec:
  privileged: false
  allowPrivilegeEscalation: false
  requiredDropCapabilities:
    - ALL
  volumes:
    - 'configMap'
    - 'emptyDir'
    - 'projected'
    - 'secret'
    - 'downwardAPI'
    - 'persistentVolumeClaim'
  runAsUser:
    rule: 'MustRunAsNonRoot'
  seLinux:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'RunAsAny'
```

## 📋 Backup Strategy

### 1. Database Backup

```bash
#!/bin/bash
# backup-db.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
DB_NAME="lunary_production"

# Create backup
kubectl exec -it postgres-pod -n lunary-prod -- \
  pg_dump -U lunary_prod -d $DB_NAME | \
  gzip > $BACKUP_DIR/lunary_backup_$DATE.sql.gz

# Retain last 7 days
find $BACKUP_DIR -name "lunary_backup_*.sql.gz" -mtime +7 -delete

# Upload to S3
aws s3 cp $BACKUP_DIR/lunary_backup_$DATE.sql.gz \
  s3://lunary-backups/database/
```

### 2. Automated Backup CronJob

```yaml
# backup-cronjob.yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: database-backup
  namespace: lunary-prod
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: postgres:15
            command:
            - /bin/sh
            - -c
            - |
              DATE=$(date +%Y%m%d_%H%M%S)
              pg_dump -h postgres-service -U $DB_USER -d lunary_production | \
              gzip > /backup/lunary_backup_$DATE.sql.gz
            env:
            - name: DB_USER
              valueFrom:
                secretKeyRef:
                  name: lunary-secrets
                  key: db-user
            - name: PGPASSWORD
              valueFrom:
                secretKeyRef:
                  name: lunary-secrets
                  key: db-password
            volumeMounts:
            - name: backup-storage
              mountPath: /backup
          volumes:
          - name: backup-storage
            persistentVolumeClaim:
              claimName: backup-pvc
          restartPolicy: OnFailure
```

## 🚀 Deployment Process

### 1. CI/CD Pipeline

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production
on:
  push:
    tags:
      - 'v*'

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Setup Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    - name: Run tests
      run: |
        cd services/auth-service
        go test -cover ./...
        bash test_health.sh
    
  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Build Docker image
      run: |
        docker build -t lunary/auth-service:${{ github.ref_name }} \
          services/auth-service/
    - name: Push to registry
      run: |
        echo ${{ secrets.DOCKER_PASSWORD }} | \
        docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
        docker push lunary/auth-service:${{ github.ref_name }}
  
  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Deploy to Kubernetes
      run: |
        echo "${{ secrets.KUBECONFIG }}" | base64 -d > kubeconfig
        export KUBECONFIG=kubeconfig
        
        # Update image tag
        sed -i "s|image: lunary/auth-service:.*|image: lunary/auth-service:${{ github.ref_name }}|" \
          deploy/k8s/auth-deployment.yaml
        
        # Apply manifests
        kubectl apply -f deploy/k8s/ -n lunary-prod
        
        # Wait for rollout
        kubectl rollout status deployment/auth-service -n lunary-prod
```

### 2. Blue-Green Deployment

```bash
#!/bin/bash
# blue-green-deploy.sh

NEW_VERSION=$1
CURRENT_VERSION=$(kubectl get deployment auth-service -n lunary-prod -o jsonpath='{.spec.template.spec.containers[0].image}' | cut -d: -f2)

echo "Deploying version $NEW_VERSION (current: $CURRENT_VERSION)"

# Deploy new version alongside current
kubectl patch deployment auth-service-green -n lunary-prod \
  -p '{"spec":{"template":{"spec":{"containers":[{"name":"auth-service","image":"lunary/auth-service:'$NEW_VERSION'"}]}}}}'

# Wait for new deployment
kubectl rollout status deployment/auth-service-green -n lunary-prod

# Health check new version
kubectl port-forward deployment/auth-service-green 8081:8081 -n lunary-prod &
sleep 5

if curl -f http://localhost:8081/health/ready; then
  echo "Health check passed, switching traffic"
  
  # Switch service to green
  kubectl patch service auth-service -n lunary-prod \
    -p '{"spec":{"selector":{"version":"green"}}}'
  
  # Scale down blue
  kubectl scale deployment auth-service-blue --replicas=0 -n lunary-prod
else
  echo "Health check failed, keeping current version"
  exit 1
fi
```

## 📈 Performance Tuning

### 1. Database Optimization

```sql
-- Optimize PostgreSQL for production
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET random_page_cost = 1.1;
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;

-- Reload configuration
SELECT pg_reload_conf();

-- Create indexes for performance
CREATE INDEX CONCURRENTLY idx_users_phone_hash ON users USING hash(phone);
CREATE INDEX CONCURRENTLY idx_feedback_created_at ON feedback(created_at DESC);
CREATE INDEX CONCURRENTLY idx_scenarios_active_title ON scenarios(is_active, title) WHERE is_active = true;
```

### 2. Application Configuration

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: auth-service-config
  namespace: lunary-prod
data:
  DB_MAX_OPEN_CONNS: "25"
  DB_MAX_IDLE_CONNS: "5"
  DB_CONN_MAX_LIFETIME: "300s"
  HTTP_READ_TIMEOUT: "15s"
  HTTP_WRITE_TIMEOUT: "15s"
  HTTP_IDLE_TIMEOUT: "60s"
  SHUTDOWN_TIMEOUT: "30s"
  LOG_LEVEL: "info"
```

## 🔍 Troubleshooting

### Common Issues

1. **Pod не стартует**
   ```bash
   kubectl describe pod <pod-name> -n lunary-prod
   kubectl logs <pod-name> -n lunary-prod
   ```

2. **Database connection issues**
   ```bash
   kubectl exec -it <auth-pod> -n lunary-prod -- \
     curl http://localhost:8081/health/ready
   ```

3. **High memory usage**
   ```bash
   kubectl top pods -n lunary-prod
   kubectl exec -it <pod-name> -n lunary-prod -- \
     curl http://localhost:8081/debug/pprof/heap
   ```

### Performance Issues

1. **Slow database queries**
   ```sql
   SELECT query, mean_time, calls 
   FROM pg_stat_statements 
   ORDER BY mean_time DESC 
   LIMIT 10;
   ```

2. **High CPU usage**
   ```bash
   kubectl exec -it <pod-name> -n lunary-prod -- \
     curl http://localhost:8081/debug/pprof/profile?seconds=30
   ```

## ✅ Production Readiness Checklist

### Pre-Launch
- [ ] **Security scan** пройден
- [ ] **Load testing** выполнен
- [ ] **Backup strategy** настроена  
- [ ] **Monitoring** настроен
- [ ] **Alerting** настроен
- [ ] **Documentation** актуальна
- [ ] **Incident response** план готов

### Post-Launch
- [ ] **Health checks** работают
- [ ] **Logs** собираются корректно
- [ ] **Metrics** отображаются в Grafana
- [ ] **Alerts** настроены и тестированы
- [ ] **Backup** выполняется автоматически
- [ ] **Performance** соответствует SLA

Этот production deployment обеспечивает **high availability**, **security**, **observability** и **scalability** для enterprise environment! 🚀
