# Deployment Guide: Ubuntu VPS Server Setup

This guide provides step-by-step instructions for deploying the Go/Gin E-Commerce API on an Ubuntu VPS server with NGINX load balancing, supporting 10,000+ concurrent users.

## Table of Contents

1. [Server Requirements](#server-requirements)
2. [Initial Server Setup](#initial-server-setup)
3. [Install Dependencies](#install-dependencies)
4. [Configure Firewall](#configure-firewall)
5. [Application Deployment](#application-deployment)
6. [NGINX Load Balancer Setup](#nginx-load-balancer-setup)
7. [SSL Certificate Setup](#ssl-certificate-setup)
8. [Database Initialization](#database-initialization)
9. [Monitoring & Maintenance](#monitoring--maintenance)
10. [Troubleshooting](#troubleshooting)

---

## Server Requirements

### Minimum Specifications
- **OS**: Ubuntu 20.04 LTS or 22.04 LTS
- **RAM**: 4GB minimum, 8GB recommended
- **CPU**: 2 cores minimum, 4 cores recommended
- **Disk**: 20GB SSD minimum
- **Network**: 100Mbps connection

### Recommended VPS Providers
- DigitalOcean (Droplet)
- AWS EC2
- Linode
- Vultr
- Google Cloud Compute Engine

---

## Initial Server Setup

### Step 1: Connect to Your Server

```bash
# SSH into your server
ssh root@your_server_ip

# Or if you have a non-root user
ssh username@your_server_ip
```

### Step 2: Update System Packages

```bash
# Update package lists
sudo apt update

# Upgrade installed packages
sudo apt upgrade -y

# Install basic utilities
sudo apt install -y curl wget git vim software-properties-common
```

### Step 3: Create a Deployment User (Optional but Recommended)

```bash
# Create new user
sudo adduser deployer

# Add to sudo group
sudo usermod -aG sudo deployer

# Switch to new user
su - deployer

# Or logout and SSH as new user
exit
ssh deployer@your_server_ip
```

### Step 4: Set Up SSH Key Authentication (Recommended)

On your local machine:
```bash
# Generate SSH key if you don't have one
ssh-keygen -t ed25519 -C "your_email@example.com"

# Copy SSH key to server
ssh-copy-id deployer@your_server_ip
```

On the server:
```bash
# Disable password authentication (after SSH key works)
sudo nano /etc/ssh/sshd_config
```

Change these lines:
```
PasswordAuthentication no
PubkeyAuthentication yes
```

Restart SSH:
```bash
sudo systemctl restart sshd
```

---

## Install Dependencies

### Step 1: Install Docker

```bash
# Remove old versions
sudo apt remove -y docker docker-engine docker.io containerd runc

# Install dependencies
sudo apt install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# Add Docker's official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Set up Docker repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io

# Start and enable Docker
sudo systemctl start docker
sudo systemctl enable docker

# Add your user to docker group
sudo usermod -aG docker $USER

# Apply group changes (or logout and login again)
newgrp docker

# Verify Docker installation
docker --version
docker run hello-world
```

### Step 2: Install Docker Compose

```bash
# Download Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

# Make it executable
sudo chmod +x /usr/local/bin/docker-compose

# Verify installation
docker-compose --version
```

### Step 3: Install NGINX

```bash
# Install NGINX
sudo apt install -y nginx

# Start and enable NGINX
sudo systemctl start nginx
sudo systemctl enable nginx

# Verify installation
nginx -v
sudo systemctl status nginx
```

---

## Configure Firewall

### Step 1: Set Up UFW (Uncomplicated Firewall)

```bash
# Allow SSH (IMPORTANT: Do this first!)
sudo ufw allow 22/tcp
sudo ufw allow OpenSSH

# Allow HTTP and HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow custom ports if needed (for development/testing)
# sudo ufw allow 8080/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status verbose
```

### Step 2: Configure Additional Security

```bash
# Install fail2ban to prevent brute force attacks
sudo apt install -y fail2ban

# Start and enable fail2ban
sudo systemctl start fail2ban
sudo systemctl enable fail2ban

# Check status
sudo systemctl status fail2ban
```

---

## Application Deployment

### Step 1: Clone Repository

```bash
# Navigate to home directory
cd ~

# Clone your repository
git clone https://github.com/your-username/gin-ecommerce-api.git

# Or if private repository, set up SSH key with GitHub first
cd gin-ecommerce-api

# Check current branch
git branch
```

### Step 2: Configure Environment Variables

```bash
# Copy example environment file
cp .env.example .env

# Edit environment file
nano .env
```

**Production .env Configuration**:

```env
# Server Configuration
SERVER_PORT=8080
ENV=production

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=ecommerce_user
DB_PASSWORD=YOUR_SECURE_PASSWORD_HERE
DB_NAME=ecommerce
DB_SSLMODE=disable

# JWT Configuration (CRITICAL: Use strong random string!)
JWT_SECRET=YOUR_VERY_SECURE_RANDOM_STRING_MINIMUM_32_CHARACTERS
JWT_EXPIRE_TIME=24

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

**Generate secure secrets**:
```bash
# Generate secure JWT secret
openssl rand -base64 32

# Generate secure database password
openssl rand -base64 24
```

### Step 3: Create Production Docker Compose File

```bash
nano docker-compose.production.yml
```

**Paste this configuration**:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: ecommerce-postgres
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    command: >
      postgres
      -c max_connections=200
      -c shared_buffers=256MB
      -c effective_cache_size=1GB
      -c maintenance_work_mem=64MB
      -c checkpoint_completion_target=0.9
      -c wal_buffers=16MB
      -c default_statistics_target=100
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/migrations:ro
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
    command: redis-server --maxmemory 512mb --maxmemory-policy allkeys-lru
    networks:
      - ecommerce-network
    restart: always
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 3

  api1:
    build: .
    container_name: ecommerce-api-1
    env_file: .env
    ports:
      - "8081:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - ecommerce-network
    restart: always
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  api2:
    build: .
    container_name: ecommerce-api-2
    env_file: .env
    ports:
      - "8082:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - ecommerce-network
    restart: always
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  api3:
    build: .
    container_name: ecommerce-api-3
    env_file: .env
    ports:
      - "8083:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - ecommerce-network
    restart: always
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  postgres_data:

networks:
  ecommerce-network:
    driver: bridge
```

### Step 4: Build and Start Services

```bash
# Build and start all services in detached mode
docker-compose -f docker-compose.production.yml up -d --build

# This will take 5-10 minutes for first build

# Check if all containers are running
docker-compose -f docker-compose.production.yml ps

# View logs
docker-compose -f docker-compose.production.yml logs -f

# Press Ctrl+C to exit log view
```

---

## NGINX Load Balancer Setup

### Step 1: Remove Default Configuration

```bash
# Backup default config (optional)
sudo cp /etc/nginx/sites-enabled/default /etc/nginx/sites-enabled/default.bak

# Remove default configuration
sudo rm /etc/nginx/sites-enabled/default
```

### Step 2: Create Load Balancer Configuration

```bash
# Create new configuration file
sudo nano /etc/nginx/sites-available/api-loadbalancer
```

**Paste this configuration** (adjust `server_name` to your domain or IP):

```nginx
upstream api_backend {
    least_conn;  # Use least connections load balancing
    server localhost:8081 max_fails=3 fail_timeout=30s;
    server localhost:8082 max_fails=3 fail_timeout=30s;
    server localhost:8083 max_fails=3 fail_timeout=30s;
    
    keepalive 32;
}

# Rate limiting zones
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=100r/s;
limit_req_zone $binary_remote_addr zone=auth_limit:10m rate=10r/s;
limit_conn_zone $binary_remote_addr zone=addr:10m;

server {
    listen 80;
    server_name your_domain.com www.your_domain.com;  # Change this!

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;

    # Connection limit
    limit_conn addr 20;

    # Logging
    access_log /var/log/nginx/api_access.log;
    error_log /var/log/nginx/api_error.log warn;

    # Health check endpoint (no rate limit)
    location /health {
        proxy_pass http://api_backend;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        access_log off;
    }

    # Authentication endpoints (stricter rate limit)
    location ~ ^/api/v1/auth/(login|register) {
        limit_req zone=auth_limit burst=20 nodelay;
        
        proxy_pass http://api_backend;
        proxy_http_version 1.1;
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
        
        # Enable gzip compression
        gzip on;
        gzip_types application/json;
        gzip_min_length 1000;
    }

    # Serve static files if needed
    location /static/ {
        alias /var/www/static/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

### Step 3: Enable Configuration

```bash
# Create symbolic link to enable site
sudo ln -s /etc/nginx/sites-available/api-loadbalancer /etc/nginx/sites-enabled/

# Test NGINX configuration
sudo nginx -t

# If test passes, reload NGINX
sudo systemctl reload nginx

# Check NGINX status
sudo systemctl status nginx
```

### Step 4: Optimize NGINX Performance

```bash
# Edit main NGINX configuration
sudo nano /etc/nginx/nginx.conf
```

Update these settings:

```nginx
user www-data;
worker_processes auto;  # Automatically use all CPU cores
pid /run/nginx.pid;

events {
    worker_connections 4096;  # Increase from default 768
    use epoll;  # Use efficient event model
    multi_accept on;
}

http {
    # Basic Settings
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    client_max_body_size 20M;

    # Rest of configuration...
}
```

Reload NGINX:
```bash
sudo nginx -t && sudo systemctl reload nginx
```

---

## SSL Certificate Setup

### Option 1: Let's Encrypt (Free, Recommended)

```bash
# Install Certbot
sudo apt install -y certbot python3-certbot-nginx

# Obtain and install SSL certificate
sudo certbot --nginx -d your_domain.com -d www.your_domain.com

# Follow the prompts:
# 1. Enter email address
# 2. Agree to terms of service
# 3. Choose to redirect HTTP to HTTPS (option 2)

# Verify auto-renewal is set up
sudo certbot renew --dry-run

# Check certbot timer
sudo systemctl status certbot.timer
```

Certbot will automatically:
- Obtain SSL certificate
- Modify NGINX configuration
- Set up auto-renewal (runs twice daily)

### Option 2: Custom SSL Certificate

If you have your own SSL certificate:

```bash
# Create SSL directory
sudo mkdir -p /etc/nginx/ssl

# Copy your certificates
sudo cp your_certificate.crt /etc/nginx/ssl/cert.pem
sudo cp your_private_key.key /etc/nginx/ssl/key.pem

# Set proper permissions
sudo chmod 600 /etc/nginx/ssl/key.pem
sudo chmod 644 /etc/nginx/ssl/cert.pem
```

Update NGINX configuration:
```bash
sudo nano /etc/nginx/sites-available/api-loadbalancer
```

Add HTTPS server block:

```nginx
server {
    listen 443 ssl http2;
    server_name your_domain.com www.your_domain.com;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Include all location blocks from port 80 config above
    # ...
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name your_domain.com www.your_domain.com;
    return 301 https://$server_name$request_uri;
}
```

---

## Database Initialization

### Step 1: Apply Database Indexes

```bash
# Navigate to project directory
cd ~/gin-ecommerce-api

# Run index migration
docker-compose -f docker-compose.production.yml exec postgres \
  psql -U ${DB_USER} -d ${DB_NAME} -f /migrations/add_indexes.sql

# Expected output: CREATE INDEX (multiple times)
```

### Step 2: Verify Database Setup

```bash
# Connect to database
docker-compose -f docker-compose.production.yml exec postgres \
  psql -U ${DB_USER} -d ${DB_NAME}

# Inside psql, run:
\dt                 # List all tables
\di                 # List all indexes
SELECT count(*) FROM users;   # Check if tables exist
\q                  # Exit psql
```

### Step 3: Create Admin User (Optional)

```bash
# Connect to database
docker-compose -f docker-compose.production.yml exec postgres \
  psql -U ${DB_USER} -d ${DB_NAME}

# Update user role to admin
UPDATE users SET role = 'admin' WHERE email = 'your_admin_email@example.com';

# Exit
\q
```

---

## Monitoring & Maintenance

### Create Monitoring Script

```bash
# Create monitoring script
nano ~/monitor.sh
```

Paste this content:

```bash
#!/bin/bash

echo "======================================"
echo "E-Commerce API Monitoring Dashboard"
echo "======================================"
echo ""

echo "=== System Resources ==="
echo "Memory Usage:"
free -h | grep -E "Mem|Swap"
echo ""
echo "Disk Usage:"
df -h / | tail -1
echo ""
echo "CPU Load:"
uptime
echo ""

echo "=== Docker Containers ==="
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
echo ""

echo "=== API Health Check ==="
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost/health)
if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ API is healthy (HTTP $HTTP_CODE)"
else
    echo "❌ API is unhealthy (HTTP $HTTP_CODE)"
fi
echo ""

echo "=== Container Resource Usage ==="
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
echo ""

echo "=== Recent NGINX Errors ==="
sudo tail -10 /var/log/nginx/api_error.log 2>/dev/null || echo "No errors found"
echo ""

echo "=== Request Rate (last minute) ==="
sudo tail -1000 /var/log/nginx/api_access.log 2>/dev/null | wc -l
echo "requests in last 1000 log entries"
```

Make it executable:
```bash
chmod +x ~/monitor.sh

# Run it
./monitor.sh
```

### Set Up Automated Monitoring

```bash
# Install monitoring tools
sudo apt install -y htop iotop ncdu

# Create cron job for daily health checks
crontab -e

# Add these lines:
# Check health every 5 minutes and log to file
*/5 * * * * curl -s http://localhost/health > /dev/null || echo "$(date): API health check failed" >> /var/log/api-health.log

# Daily backup at 2 AM
0 2 * * * cd ~/gin-ecommerce-api && docker-compose -f docker-compose.production.yml exec -T postgres pg_dump -U ${DB_USER} ${DB_NAME} > ~/backups/db_backup_$(date +\%Y\%m\%d).sql
```

### Log Rotation

```bash
# Create logrotate configuration
sudo nano /etc/logrotate.d/api-logs
```

Add:
```
/var/log/nginx/api_*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 www-data adm
    sharedscripts
    postrotate
        [ -f /var/run/nginx.pid ] && kill -USR1 `cat /var/run/nginx.pid`
    endscript
}
```

---

## Troubleshooting

### Common Issues

#### 1. Containers Won't Start

```bash
# Check logs for specific container
docker-compose -f docker-compose.production.yml logs api1

# Check if ports are already in use
sudo netstat -tulpn | grep -E '8081|8082|8083'

# Restart services
docker-compose -f docker-compose.production.yml restart
```

#### 2. NGINX Configuration Errors

```bash
# Test configuration
sudo nginx -t

# View error log
sudo tail -50 /var/log/nginx/error.log

# Restart NGINX
sudo systemctl restart nginx
```

#### 3. Database Connection Issues

```bash
# Check if PostgreSQL is running
docker-compose -f docker-compose.production.yml exec postgres pg_isready -U ${DB_USER}

# Check container logs
docker-compose -f docker-compose.production.yml logs postgres

# Restart database
docker-compose -f docker-compose.production.yml restart postgres
```

#### 4. 502 Bad Gateway

```bash
# Check if backend services are running
docker ps

# Check backend logs
docker-compose -f docker-compose.production.yml logs api1 api2 api3

# Verify upstream configuration
sudo nginx -T | grep upstream

# Test backend directly
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

#### 5. High Memory Usage

```bash
# Check memory usage
free -h
docker stats

# Clear Docker cache
docker system prune -a

# Restart specific container
docker-compose -f docker-compose.production.yml restart api1
```

### Useful Commands

```bash
# View all container logs
docker-compose -f docker-compose.production.yml logs -f

# Restart all services
docker-compose -f docker-compose.production.yml restart

# Stop all services
docker-compose -f docker-compose.production.yml down

# Rebuild and restart
docker-compose -f docker-compose.production.yml up -d --build

# Check disk usage
docker system df

# Clean up Docker resources
docker system prune -a --volumes

# Export database backup
docker-compose -f docker-compose.production.yml exec postgres \
  pg_dump -U ${DB_USER} ${DB_NAME} > backup_$(date +%Y%m%d).sql

# Restore database backup
docker-compose -f docker-compose.production.yml exec -T postgres \
  psql -U ${DB_USER} -d ${DB_NAME} < backup_20240101.sql
```

---

## Update/Deployment Process

### Deploy New Version

```bash
# Navigate to project directory
cd ~/gin-ecommerce-api

# Pull latest code
git pull origin main

# Rebuild and restart services (zero-downtime)
docker-compose -f docker-compose.production.yml up -d --build --no-deps api1
sleep 30
docker-compose -f docker-compose.production.yml up -d --build --no-deps api2
sleep 30
docker-compose -f docker-compose.production.yml up -d --build --no-deps api3

# Verify all services are healthy
docker-compose -f docker-compose.production.yml ps
```

### Rollback Process

```bash
# View commit history
git log --oneline -10

# Rollback to previous version
git reset --hard <commit-hash>

# Rebuild services
docker-compose -f docker-compose.production.yml up -d --build
```

---

## Performance Optimization

### System Tuning

```bash
# Edit sysctl configuration
sudo nano /etc/sysctl.conf
```

Add these optimizations:

```conf
# Increase maximum connections
net.core.somaxconn = 65535

# Increase file descriptors
fs.file-max = 65535

# TCP optimization
net.ipv4.tcp_max_syn_backlog = 8192
net.ipv4.tcp_fin_timeout = 30
net.ipv4.tcp_keepalive_time = 300

# Enable BBR congestion control
net.core.default_qdisc = fq
net.ipv4.tcp_congestion_control = bbr
```

Apply changes:
```bash
sudo sysctl -p
```

### Increase File Limits

```bash
# Edit limits configuration
sudo nano /etc/security/limits.conf
```

Add:
```
* soft nofile 65535
* hard nofile 65535
* soft nproc 65535
* hard nproc 65535
```

Logout and login again for changes to take effect.

---

## Security Checklist

- [ ] SSH key authentication enabled
- [ ] Password authentication disabled
- [ ] Firewall (UFW) configured and enabled
- [ ] Fail2ban installed and running
- [ ] SSL certificate installed and auto-renewal working
- [ ] Strong passwords for database and JWT secret
- [ ] Environment variables properly secured (not in git)
- [ ] Regular backups configured
- [ ] Log rotation configured
- [ ] Monitoring and alerts set up
- [ ] Keep system and Docker images updated

---

## Post-Deployment Verification

```bash
# 1. Check all services are running
docker ps

# 2. Test API health
curl http://your_domain.com/health

# 3. Test load balancer distribution
for i in {1..10}; do curl -s http://your_domain.com/health; done

# 4. Check SSL certificate
curl -I https://your_domain.com

# 5. Test API endpoints
curl http://your_domain.com/api/v1/products

# 6. Monitor logs
sudo tail -f /var/log/nginx/api_access.log

# 7. Run performance test
ab -n 1000 -c 100 http://your_domain.com/api/v1/products
```

---

## Next Steps

1. Set up automated backups to cloud storage (S3, Google Cloud Storage)
2. Implement monitoring with Prometheus + Grafana
3. Set up log aggregation with ELK stack or similar
4. Configure CDN for static assets (Cloudflare, AWS CloudFront)
5. Set up CI/CD pipeline for automated deployments
6. Implement blue-green or canary deployment strategy
7. Set up database replication for high availability

---

## Support Resources

- **Application Docs**: See README.md, TUTORIAL.md, SCALING.md
- **Docker Docs**: https://docs.docker.com/
- **NGINX Docs**: https://nginx.org/en/docs/
- **PostgreSQL Docs**: https://www.postgresql.org/docs/
- **Ubuntu Server Guide**: https://ubuntu.com/server/docs

For application-specific issues, refer to TROUBLESHOOTING section or check application logs.

---

**Deployment Checklist Complete!** 🚀

Your Go/Gin E-Commerce API is now deployed and ready to handle 10,000+ concurrent users from your Next.js frontend!
