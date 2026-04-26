# Scaling Guide: Supporting 10,000+ Concurrent Users

This guide provides strategies and implementations to scale the e-commerce API to handle high traffic from a Next.js frontend.

## Table of Contents

1. [Current Bottlenecks](#current-bottlenecks)
2. [Scaling Strategy Overview](#scaling-strategy-overview)
3. [Database Optimization](#database-optimization)
4. [Caching Layer](#caching-layer)
5. [Load Balancing](#load-balancing)
6. [Horizontal Scaling](#horizontal-scaling)
7. [Performance Optimizations](#performance-optimizations)
8. [Monitoring & Observability](#monitoring--observability)
9. [Infrastructure Setup](#infrastructure-setup)

---

## Current Bottlenecks

Before scaling, identify potential bottlenecks:

1. **Database Connections**: Single PostgreSQL instance with limited connections
2. **Single Instance**: No horizontal scaling capability
3. **No Caching**: Every request hits the database
4. **JWT Validation**: Token validation on every request
5. **Stock Management**: Race conditions under high load
6. **No Rate Limiting**: Vulnerable to abuse
7. **No CDN**: Static content served from application

---

## Scaling Strategy Overview

### Architecture Layers

```
┌─────────────────────────────────────────────────────┐
│                    CDN (Cloudflare)                  │
│              Static Assets & API Gateway             │
└─────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────┐
│              Load Balancer (NGINX)                   │
│           SSL Termination & Rate Limiting            │
└─────────────────────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        ▼                 ▼                 ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  API Node 1  │  │  API Node 2  │  │  API Node N  │
│   (Go/Gin)   │  │   (Go/Gin)   │  │   (Go/Gin)   │
└──────────────┘  └──────────────┘  └──────────────┘
        │                 │                 │
        └─────────────────┼─────────────────┘
                          ▼
        ┌─────────────────────────────────────┐
        │      Redis Cluster (Caching)         │
        │   Session Store & Cache Layer        │
        └─────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────────┐
        │   PostgreSQL (Primary + Replicas)    │
        │      Read Replicas for Scaling       │
        └─────────────────────────────────────┘
```

### Key Improvements

1. **Horizontal Scaling**: Multiple API instances behind load balancer
2. **Caching**: Redis for session management and data caching
3. **Database**: Read replicas for read-heavy operations
4. **Connection Pooling**: Optimized database connections
5. **Rate Limiting**: Protect against abuse
6. **CDN**: Offload static assets
7. **Monitoring**: Real-time metrics and alerts

---

## Database Optimization

### 1. Connection Pooling

Update `internal/database/database.go`:

```go
func Connect(cfg *config.Config) error {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.DBName,
        cfg.Database.SSLMode,
    )

    logLevel := logger.Info
    if cfg.IsProduction() {
        logLevel = logger.Error
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logLevel),
        PrepareStmt: true, // Enable prepared statement cache
    })
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }

    // Connection pool settings for high load
    sqlDB.SetMaxIdleConns(25)           // Idle connections
    sqlDB.SetMaxOpenConns(100)          // Max open connections
    sqlDB.SetConnMaxLifetime(time.Hour) // Connection lifetime
    sqlDB.SetConnMaxIdleTime(10 * time.Minute)

    DB = db
    log.Println("Database connection established with optimized pool")
    return nil
}
```

### 2. Database Indexes

Create migration file `migrations/add_indexes.sql`:

```sql
-- Users table indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- Products table indexes
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);
CREATE INDEX IF NOT EXISTS idx_products_price ON products(price);
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at DESC);

-- Orders table indexes
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);

-- Cart items indexes
CREATE INDEX IF NOT EXISTS idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_product_id ON cart_items(product_id);

-- Order items indexes
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);

-- Posts indexes
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_is_active ON posts(is_active);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_products_category_active ON products(category, is_active);
CREATE INDEX IF NOT EXISTS idx_orders_user_status ON orders(user_id, status);
```

### 3. Read Replicas Configuration

Add to `.env`:

```env
# Database Read Replicas
DB_READ_REPLICA_HOSTS=replica1.example.com,replica2.example.com
```

---

## Caching Layer

### 1. Redis Integration

Add Redis configuration to `internal/config/config.go`:

```go
type RedisConfig struct {
    Host     string
    Port     string
    Password string
    DB       int
}

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
    Redis    RedisConfig // Add this
}

func Load() *Config {
    return &Config{
        // ... existing config
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnv("REDIS_PORT", "6379"),
            Password: getEnv("REDIS_PASSWORD", ""),
            DB:       getEnvAsInt("REDIS_DB", 0),
        },
    }
}
```

### 2. Redis Client Setup

Create `internal/cache/redis.go`:

```go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/dat19/gin-ecommerce-api/internal/config"
    "github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect(cfg *config.Config) error {
    Client = redis.NewClient(&redis.Options{
        Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
        Password:     cfg.Redis.Password,
        DB:           cfg.Redis.DB,
        PoolSize:     50,
        MinIdleConns: 10,
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := Client.Ping(ctx).Err(); err != nil {
        return fmt.Errorf("failed to connect to Redis: %w", err)
    }

    return nil
}

func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return Client.Set(ctx, key, data, ttl).Err()
}

func Get(ctx context.Context, key string, dest interface{}) error {
    data, err := Client.Get(ctx, key).Bytes()
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dest)
}

func Delete(ctx context.Context, keys ...string) error {
    return Client.Del(ctx, keys...).Err()
}

func Close() error {
    if Client != nil {
        return Client.Close()
    }
    return nil
}
```

### 3. Cache Middleware

Create `internal/api/middleware/cache.go`:

```go
package middleware

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"

    "github.com/dat19/gin-ecommerce-api/internal/cache"
    "github.com/gin-gonic/gin"
)

// CacheResponse caches GET request responses
func CacheResponse(ttl time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Only cache GET requests
        if c.Request.Method != "GET" {
            c.Next()
            return
        }

        // Generate cache key
        cacheKey := generateCacheKey(c)

        // Try to get from cache
        var cachedResponse map[string]interface{}
        err := cache.Get(c.Request.Context(), cacheKey, &cachedResponse)
        if err == nil {
            c.JSON(200, cachedResponse)
            c.Abort()
            return
        }

        // Create response writer wrapper to capture response
        blw := &bodyLogWriter{body: []byte{}, ResponseWriter: c.Writer}
        c.Writer = blw

        c.Next()

        // Cache successful responses
        if c.Writer.Status() == 200 {
            go func() {
                ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
                defer cancel()
                cache.Set(ctx, cacheKey, blw.body, ttl)
            }()
        }
    }
}

func generateCacheKey(c *gin.Context) string {
    url := c.Request.URL.String()
    hash := sha256.Sum256([]byte(url))
    return fmt.Sprintf("cache:%s", hex.EncodeToString(hash[:]))
}

type bodyLogWriter struct {
    gin.ResponseWriter
    body []byte
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
    w.body = append(w.body, b...)
    return w.ResponseWriter.Write(b)
}
```

---

## Load Balancing

### NGINX Configuration

Create `nginx/nginx.conf`:

```nginx
upstream api_backend {
    least_conn;  # Load balancing method
    server api1:8080 max_fails=3 fail_timeout=30s;
    server api2:8080 max_fails=3 fail_timeout=30s;
    server api3:8080 max_fails=3 fail_timeout=30s;
    
    keepalive 32;
}

# Rate limiting zones
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=100r/s;
limit_req_zone $binary_remote_addr zone=auth_limit:10m rate=10r/s;

server {
    listen 80;
    server_name api.yourdomain.com;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Health check endpoint (no rate limit)
    location /health {
        proxy_pass http://api_backend;
        access_log off;
    }

    # Authentication endpoints (stricter rate limit)
    location ~ ^/api/v1/auth/(login|register) {
        limit_req zone=auth_limit burst=20 nodelay;
        
        proxy_pass http://api_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # General API endpoints
    location /api/ {
        limit_req zone=api_limit burst=200 nodelay;
        
        proxy_pass http://api_backend;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # Connection keep-alive
        proxy_buffering off;
    }
}
```

---

## Horizontal Scaling

### Docker Compose for Production with Multiple Instances

Create `docker-compose.scaled.yml`:

```yaml
version: '3.8'

services:
  nginx:
    image: nginx:alpine
    container_name: ecommerce-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
    depends_on:
      - api1
      - api2
      - api3
    networks:
      - ecommerce-network
    restart: always

  postgres-primary:
    image: postgres:15-alpine
    container_name: ecommerce-postgres-primary
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_MAX_CONNECTIONS: 200
    command: >
      postgres
      -c max_connections=200
      -c shared_buffers=256MB
      -c effective_cache_size=1GB
      -c maintenance_work_mem=64MB
      -c checkpoint_completion_target=0.9
      -c wal_buffers=16MB
      -c default_statistics_target=100
      -c random_page_cost=1.1
      -c effective_io_concurrency=200
      -c work_mem=2621kB
      -c min_wal_size=1GB
      -c max_wal_size=4GB
    volumes:
      - postgres_primary_data:/var/lib/postgresql/data
    networks:
      - ecommerce-network
    restart: always
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: ecommerce-redis
    command: >
      redis-server
      --maxmemory 512mb
      --maxmemory-policy allkeys-lru
      --save ""
      --appendonly no
    ports:
      - "6379:6379"
    networks:
      - ecommerce-network
    restart: always
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 3

  api1:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ecommerce-api-1
    environment:
      ENV: production
      SERVER_PORT: 8080
      DB_HOST: postgres-primary
      DB_PORT: 5432
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: ${DB_NAME}
      DB_SSLMODE: disable
      JWT_SECRET: ${JWT_SECRET}
      JWT_EXPIRE_TIME: ${JWT_EXPIRE_TIME:-24}
      REDIS_HOST: redis
      REDIS_PORT: 6379
    depends_on:
      postgres-primary:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - ecommerce-network
    restart: always
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M

  api2:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ecommerce-api-2
    environment:
      ENV: production
      SERVER_PORT: 8080
      DB_HOST: postgres-primary
      DB_PORT: 5432
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: ${DB_NAME}
      DB_SSLMODE: disable
      JWT_SECRET: ${JWT_SECRET}
      JWT_EXPIRE_TIME: ${JWT_EXPIRE_TIME:-24}
      REDIS_HOST: redis
      REDIS_PORT: 6379
    depends_on:
      postgres-primary:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - ecommerce-network
    restart: always
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M

  api3:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ecommerce-api-3
    environment:
      ENV: production
      SERVER_PORT: 8080
      DB_HOST: postgres-primary
      DB_PORT: 5432
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: ${DB_NAME}
      DB_SSLMODE: disable
      JWT_SECRET: ${JWT_SECRET}
      JWT_EXPIRE_TIME: ${JWT_EXPIRE_TIME:-24}
      REDIS_HOST: redis
      REDIS_PORT: 6379
    depends_on:
      postgres-primary:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - ecommerce-network
    restart: always
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M

volumes:
  postgres_primary_data:

networks:
  ecommerce-network:
    driver: bridge
```

---

## Performance Optimizations

### 1. Rate Limiting Middleware

Create `internal/api/middleware/ratelimit.go`:

```go
package middleware

import (
    "net/http"
    "sync"
    "time"

    "github.com/dat19/gin-ecommerce-api/pkg/utils"
    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

type visitor struct {
    limiter  *rate.Limiter
    lastSeen time.Time
}

var (
    visitors = make(map[string]*visitor)
    mu       sync.RWMutex
)

// RateLimit creates a rate limiting middleware
func RateLimit(requestsPerSecond int, burst int) gin.HandlerFunc {
    // Clean up old visitors every minute
    go cleanupVisitors()

    return func(c *gin.Context) {
        ip := c.ClientIP()
        
        mu.Lock()
        v, exists := visitors[ip]
        if !exists {
            limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burst)
            visitors[ip] = &visitor{limiter, time.Now()}
            v = visitors[ip]
        }
        v.lastSeen = time.Now()
        mu.Unlock()

        if !v.limiter.Allow() {
            utils.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded")
            c.Abort()
            return
        }

        c.Next()
    }
}

func cleanupVisitors() {
    for {
        time.Sleep(time.Minute)
        mu.Lock()
        for ip, v := range visitors {
            if time.Since(v.lastSeen) > 3*time.Minute {
                delete(visitors, ip)
            }
        }
        mu.Unlock()
    }
}
```

### 2. Optimize Product Listing

Update `internal/api/handlers/product.go`:

```go
func (h *ProductHandler) GetAll(c *gin.Context) {
    // Try cache first
    cacheKey := fmt.Sprintf("products:%s", c.Request.URL.Query().Encode())
    
    var products []models.Product
    err := cache.Get(c.Request.Context(), cacheKey, &products)
    if err == nil {
        utils.SuccessResponse(c, http.StatusOK, "Products retrieved from cache", products)
        return
    }

    // Query database with pagination
    page := c.DefaultQuery("page", "1")
    limit := c.DefaultQuery("limit", "20")
    
    query := database.DB.Where("is_active = ?", true)

    if category := c.Query("category"); category != "" {
        query = query.Where("category = ?", category)
    }

    // Use offset pagination
    var offset int
    fmt.Sscanf(page, "%d", &offset)
    var limitInt int
    fmt.Sscanf(limit, "%d", &limitInt)
    offset = (offset - 1) * limitInt

    if err := query.Limit(limitInt).Offset(offset).Find(&products).Error; err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve products")
        return
    }

    // Cache for 5 minutes
    go cache.Set(context.Background(), cacheKey, products, 5*time.Minute)

    utils.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", products)
}
```

---

## Monitoring & Observability

### Prometheus Metrics

Create `internal/metrics/metrics.go`:

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    RequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "api_requests_total",
            Help: "Total number of API requests",
        },
        []string{"method", "endpoint", "status"},
    )

    RequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "api_request_duration_seconds",
            Help:    "API request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    DatabaseQueries = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "database_queries_total",
            Help: "Total number of database queries",
        },
        []string{"operation"},
    )

    CacheHits = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "cache_hits_total",
            Help: "Total number of cache hits",
        },
    )

    CacheMisses = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "cache_misses_total",
            Help: "Total number of cache misses",
        },
    )
)
```

Add metrics endpoint to `cmd/main.go`:

```go
import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// In main() after routes setup
r.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

---

## Infrastructure Setup

### Deployment Steps

1. **Update Dependencies**
```bash
go get github.com/redis/go-redis/v9
go get github.com/prometheus/client_golang/prometheus
go get golang.org/x/time/rate
go mod tidy
```

2. **Create Environment File**
```bash
cp .env.example .env.production
# Edit with production values
```

3. **Build and Deploy**
```bash
# Start scaled infrastructure
docker-compose -f docker-compose.scaled.yml up -d --build

# Verify all services
docker-compose -f docker-compose.scaled.yml ps

# Check logs
docker-compose -f docker-compose.scaled.yml logs -f
```

4. **Run Database Migrations**
```bash
docker-compose -f docker-compose.scaled.yml exec postgres-primary psql -U postgres -d ecommerce -f /migrations/add_indexes.sql
```

### Load Testing

```bash
# Install Apache Bench
# Test API performance
ab -n 10000 -c 100 http://localhost/api/v1/products

# Install hey for better testing
go install github.com/rakyll/hey@latest
hey -n 50000 -c 500 http://localhost/api/v1/products
```

---

## Expected Performance

With these optimizations, you should achieve:

- **Throughput**: 5,000-10,000 requests/second
- **Response Time**: < 50ms (p95)
- **Concurrent Users**: 10,000+ simultaneous connections
- **Database Load**: Distributed across read replicas
- **Cache Hit Rate**: 70-90% for product listings
- **Uptime**: 99.9% with health checks and auto-restart

## Next Steps

1. Implement all database indexes
2. Add Redis caching layer
3. Set up NGINX load balancer
4. Deploy multiple API instances
5. Configure monitoring with Prometheus/Grafana
6. Set up auto-scaling based on metrics
7. Add CDN for static assets
8. Implement circuit breakers for fault tolerance

For production deployment, consider using Kubernetes for better orchestration and auto-scaling capabilities.
