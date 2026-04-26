# Tutorial: Getting Started with the E-Commerce API

This tutorial will guide you through setting up, running, and using the Go/Gin E-Commerce API from scratch.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Running the Application](#running-the-application)
4. [Understanding the API](#understanding-the-api)
5. [Step-by-Step Usage Guide](#step-by-step-usage-guide)
6. [Testing](#testing)
7. [Troubleshooting](#troubleshooting)

---

## Prerequisites

Before you begin, ensure you have the following installed on your system:

### Required Software

1. **Go** (version 1.21 or higher)
   - Download from: https://golang.org/dl/
   - Verify installation: `go version`

2. **Docker & Docker Compose** (for containerized deployment)
   - Download from: https://www.docker.com/products/docker-desktop
   - Verify installation: `docker --version` and `docker-compose --version`

3. **PostgreSQL** (if running without Docker)
   - Download from: https://www.postgresql.org/download/
   - Version 15 or higher recommended

4. **Git** (for cloning the repository)
   - Download from: https://git-scm.com/downloads

### Optional Tools

- **Postman** or **cURL** for testing API endpoints
- **pgAdmin** for database management
- **VS Code** with Go extension for development

---

## Installation

### Step 1: Clone the Repository

```bash
git clone <your-repository-url>
cd gin-ecommerce-api
```

### Step 2: Install Go Dependencies

```bash
go mod download
```

This will download all required packages including:
- Gin web framework
- GORM (ORM for database)
- JWT authentication
- PostgreSQL driver
- Bcrypt for password hashing

### Step 3: Configure Environment Variables

Copy the example environment file:

```bash
# On Windows (PowerShell)
Copy-Item .env.example .env

# On macOS/Linux
cp .env.example .env
```

Edit the `.env` file with your configuration:

```env
# Server Configuration
SERVER_PORT=8080
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=ecommerce
DB_SSLMODE=disable

# JWT Configuration (IMPORTANT: Change in production!)
JWT_SECRET=change-this-to-a-secure-random-string
JWT_EXPIRE_TIME=24
```

**Security Note:** Always use a strong, random JWT secret in production!

---

## Running the Application

You have three options for running the application:

### Option 1: Run with Docker (Recommended for Beginners)

This is the easiest way to get started as it includes the database.

```bash
docker-compose -f docker-compose.dev.yml up --build
```

The API will be available at: `http://localhost:8080`

To stop the application:
```bash
# Press Ctrl+C, then run:
docker-compose -f docker-compose.dev.yml down
```

### Option 2: Run Locally (Without Docker)

**Step 1:** Start PostgreSQL

Make sure PostgreSQL is running on your machine. Create a database:

```sql
CREATE DATABASE ecommerce;
```

**Step 2:** Run the application

```bash
go run cmd/main.go
```

The application will:
1. Connect to the database
2. Automatically create all necessary tables
3. Start the server on port 8080

### Option 3: Run with Go Build

Build the application first:

```bash
go build -o ecommerce-api cmd/main.go

# On Windows
.\ecommerce-api.exe

# On macOS/Linux
./ecommerce-api
```

---

## Understanding the API

### Architecture Overview

```
gin-ecommerce-api/
├── cmd/                    # Application entry point
│   └── main.go            # Server initialization
├── internal/              # Private application code
│   ├── api/
│   │   ├── handlers/      # HTTP request handlers
│   │   ├── middleware/    # Authentication, CORS, logging
│   │   └── routes/        # Route definitions
│   ├── config/            # Configuration management
│   ├── database/          # Database connection & migrations
│   └── models/            # Data models (User, Product, etc.)
└── pkg/                   # Reusable utilities
    └── utils/             # JWT, password hashing, responses
```

### Database Schema

The application manages the following entities:

- **Users**: User accounts with authentication
- **Products**: Product catalog with stock management
- **Posts**: User-generated content
- **Carts**: Shopping carts (one per user)
- **Cart Items**: Products in shopping carts
- **Orders**: Customer orders
- **Order Items**: Products in orders

All tables are created automatically when the application starts.

### API Endpoints Overview

| Category | Endpoint | Method | Auth Required |
|----------|----------|--------|---------------|
| Health | `/health` | GET | No |
| **Authentication** |
| Register | `/api/v1/auth/register` | POST | No |
| Login | `/api/v1/auth/login` | POST | No |
| Get Current User | `/api/v1/auth/me` | GET | Yes |
| Logout | `/api/v1/auth/logout` | POST | Yes |
| **Products** |
| List Products | `/api/v1/products` | GET | No |
| Get Product | `/api/v1/products/:id` | GET | No |
| Create Product | `/api/v1/products` | POST | Admin |
| Update Product | `/api/v1/products/:id` | PUT | Admin |
| Delete Product | `/api/v1/products/:id` | DELETE | Admin |
| **Cart** |
| Get Cart | `/api/v1/cart` | GET | Yes |
| Add to Cart | `/api/v1/cart/items` | POST | Yes |
| Update Cart Item | `/api/v1/cart/items/:itemId` | PUT | Yes |
| Remove from Cart | `/api/v1/cart/items/:itemId` | DELETE | Yes |
| Clear Cart | `/api/v1/cart` | DELETE | Yes |
| **Orders** |
| Create Order | `/api/v1/orders` | POST | Yes |
| List Orders | `/api/v1/orders` | GET | Yes |
| Get Order | `/api/v1/orders/:id` | GET | Yes |
| Cancel Order | `/api/v1/orders/:id/cancel` | POST | Yes |
| Update Status | `/api/v1/orders/:id/status` | PUT | Admin |
| **Posts** |
| List Posts | `/api/v1/posts` | GET | No |
| Get Post | `/api/v1/posts/:id` | GET | No |
| Create Post | `/api/v1/posts` | POST | Yes |
| Update Post | `/api/v1/posts/:id` | PUT | Yes (Owner) |
| Delete Post | `/api/v1/posts/:id` | DELETE | Yes (Owner) |
| **Users** |
| List Users | `/api/v1/users` | GET | Yes |
| Get User | `/api/v1/users/:id` | GET | Yes |
| Update User | `/api/v1/users/:id` | PUT | Yes (Self) |
| Delete User | `/api/v1/users/:id` | DELETE | Admin |

---

## Step-by-Step Usage Guide

### Step 1: Verify the Server is Running

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok"
}
```

### Step 2: Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "username": "johndoe",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

Expected response:
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "john@example.com",
      "username": "johndoe",
      "first_name": "John",
      "last_name": "Doe",
      "role": "user",
      "is_active": true
    }
  }
}
```

**Save the token!** You'll need it for authenticated requests.

### Step 3: Login (If You Already Have an Account)

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Step 4: Create Products (Admin Only)

First, you need to manually set your user role to "admin" in the database:

```sql
UPDATE users SET role = 'admin' WHERE email = 'john@example.com';
```

Then create products:

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "stock": 50,
    "category": "electronics",
    "image_url": "https://example.com/laptop.jpg"
  }'
```

### Step 5: Browse Products

```bash
# Get all products
curl http://localhost:8080/api/v1/products

# Filter by category
curl http://localhost:8080/api/v1/products?category=electronics

# Get specific product
curl http://localhost:8080/api/v1/products/1
```

### Step 6: Add Products to Cart

```bash
curl -X POST http://localhost:8080/api/v1/cart/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
```

### Step 7: View Your Cart

```bash
curl http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

Expected response:
```json
{
  "success": true,
  "message": "Cart retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "cart_items": [
      {
        "id": 1,
        "product_id": 1,
        "name": "Laptop",
        "price": 999.99,
        "quantity": 2,
        "subtotal": 1999.98
      }
    ],
    "total_price": 1999.98
  }
}
```

### Step 8: Update Cart Item Quantity

```bash
curl -X PUT http://localhost:8080/api/v1/cart/items/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "quantity": 3
  }'
```

### Step 9: Create an Order

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "shipping_address": "123 Main St, New York, NY 10001",
    "payment_method": "credit_card"
  }'
```

This will:
1. Create an order from your cart items
2. Reduce product stock
3. Clear your cart

### Step 10: View Your Orders

```bash
# Get all your orders
curl http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Get specific order
curl http://localhost:8080/api/v1/orders/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Step 11: Create a Post

```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my post."
  }'
```

### Step 12: Cancel an Order

```bash
curl -X POST http://localhost:8080/api/v1/orders/1/cancel \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

This will restore the product stock.

---

## Testing

### Running Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests for specific package
go test ./pkg/utils/... -v

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Manual API Testing with Postman

1. Import the API endpoints into Postman
2. Create an environment variable for your token
3. Use `{{token}}` in the Authorization header

**Example Postman Setup:**

```
Authorization: Bearer {{token}}
```

### Testing with Different Environments

```bash
# Development
docker-compose -f docker-compose.dev.yml up

# Staging
docker-compose -f docker-compose.staging.yml up

# Production
docker-compose -f docker-compose.prod.yml up
```

---

## Troubleshooting

### Common Issues

#### 1. Database Connection Error

**Error:** `Failed to connect to database`

**Solution:**
- Verify PostgreSQL is running
- Check database credentials in `.env`
- Ensure database exists: `CREATE DATABASE ecommerce;`

#### 2. Port Already in Use

**Error:** `bind: address already in use`

**Solution:**
```bash
# Find process using port 8080
# Windows (PowerShell)
Get-Process -Id (Get-NetTCPConnection -LocalPort 8080).OwningProcess

# macOS/Linux
lsof -i :8080

# Kill the process or change SERVER_PORT in .env
```

#### 3. JWT Token Invalid

**Error:** `Invalid or expired token`

**Solution:**
- Token expired (default 24 hours)
- Login again to get a new token
- Verify you're using `Bearer` prefix: `Bearer YOUR_TOKEN`

#### 4. Module Not Found

**Error:** `cannot find module`

**Solution:**
```bash
go mod tidy
go mod download
```

#### 5. Docker Build Fails

**Solution:**
```bash
# Clean up Docker
docker-compose down -v
docker system prune -a

# Rebuild
docker-compose -f docker-compose.dev.yml up --build
```

### Debugging Tips

1. **Enable verbose logging:**
   Set `ENV=development` in `.env`

2. **Check application logs:**
   ```bash
   docker-compose -f docker-compose.dev.yml logs -f app
   ```

3. **Inspect database:**
   ```bash
   # Connect to PostgreSQL in Docker
   docker-compose -f docker-compose.dev.yml exec postgres psql -U postgres -d ecommerce
   
   # List tables
   \dt
   
   # View users
   SELECT * FROM users;
   ```

4. **Test individual endpoints:**
   Use curl with verbose output:
   ```bash
   curl -v http://localhost:8080/api/v1/products
   ```

### Getting Help

- Check the main [README.md](README.md) for detailed API documentation
- Review [AGENTS.md](AGENTS.md) for architecture details
- Look at test files for usage examples
- Check GitHub issues or create a new one

---

## Next Steps

Now that you have the application running, you can:

1. **Explore the Code:**
   - Review handler implementations in `internal/api/handlers/`
   - Study the middleware in `internal/api/middleware/`
   - Examine models in `internal/models/`

2. **Customize the Application:**
   - Add new product categories
   - Implement product reviews
   - Add image upload functionality
   - Implement payment gateway integration

3. **Deploy to Production:**
   - Set up a production PostgreSQL database
   - Configure environment variables for production
   - Use `docker-compose.prod.yml` for deployment
   - Set up reverse proxy (nginx) and SSL/TLS

4. **Add More Features:**
   - Email notifications for orders
   - Password reset functionality
   - Product search with filters
   - Order tracking
   - Admin dashboard

---

## Useful Commands Reference

```bash
# Development
go run cmd/main.go                    # Run application
go test ./...                         # Run all tests
go build -o app cmd/main.go          # Build binary
go fmt ./...                          # Format code
go mod tidy                           # Clean dependencies

# Docker
docker-compose -f docker-compose.dev.yml up --build    # Start with rebuild
docker-compose -f docker-compose.dev.yml down          # Stop containers
docker-compose -f docker-compose.dev.yml down -v       # Stop and remove volumes
docker-compose -f docker-compose.dev.yml logs -f app   # View logs

# Database
psql -U postgres -d ecommerce         # Connect to database
\dt                                    # List tables
\d users                              # Describe users table
SELECT * FROM users;                  # Query users

# Git
git status                            # Check status
git add .                             # Stage changes
git commit -m "message"               # Commit changes
git push origin main                  # Push to remote
```

---

## Conclusion

Congratulations! You now have a fully functional e-commerce API running. This tutorial covered:

✅ Installation and setup
✅ Running the application in different modes
✅ Understanding the API structure
✅ Complete usage workflow from registration to order creation
✅ Testing and troubleshooting

For more detailed information, refer to:
- [README.md](README.md) - Complete API reference
- [AGENTS.md](AGENTS.md) - Architecture and development guide

Happy coding! 🚀
