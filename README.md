# E-Commerce API with Go/Gin

A scalable, production-ready e-commerce API built with Go and the Gin framework. Features include user authentication, product management, shopping cart, order processing, and more.

## Features

- **Authentication & Authorization**
  - User registration and login with JWT tokens
  - Role-based access control (user/admin)
  - Secure password hashing with bcrypt

- **User Management**
  - CRUD operations for users
  - User profiles with personal information

- **Product Management**
  - Full CRUD operations for products
  - Category filtering
  - Stock management
  - Admin-only product management

- **Posts**
  - Create, read, update, delete posts
  - User-owned content

- **Shopping Cart**
  - Add/remove/update items
  - Automatic cart creation on registration
  - Stock validation

- **Order Management**
  - Create orders from cart
  - Order status tracking
  - Order cancellation with stock restoration
  - Admin order management

- **Middleware**
  - Authentication middleware
  - CORS support
  - Custom logging
  - Admin authorization

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT
- **Containerization**: Docker & Docker Compose

## Project Structure

```
gin-ecommerce-api/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/           # Request handlers
│   │   ├── middleware/         # Middleware functions
│   │   └── routes/             # Route definitions
│   ├── config/                 # Configuration management
│   ├── database/               # Database connection & migrations
│   ├── models/                 # Data models
│   ├── repository/             # Data access layer (optional)
│   └── service/                # Business logic layer (optional)
├── pkg/
│   ├── utils/                  # Utility functions (JWT, password, etc.)
│   └── validator/              # Custom validators
├── docker-compose.dev.yml      # Development environment
├── docker-compose.staging.yml  # Staging environment
├── docker-compose.prod.yml     # Production environment
├── Dockerfile                  # Multi-stage Docker build
├── .env.example                # Environment variables template
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose
- PostgreSQL (if running locally without Docker)

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd gin-ecommerce-api
```

2. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Install dependencies**
```bash
go mod download
```

### Running Locally

**Option 1: With Docker Compose (Recommended)**

Development environment:
```bash
docker-compose -f docker-compose.dev.yml up --build
```

Staging environment:
```bash
docker-compose -f docker-compose.staging.yml up --build
```

Production environment:
```bash
docker-compose -f docker-compose.prod.yml up --build
```

**Option 2: Without Docker**

1. Start PostgreSQL locally
2. Update `.env` with your database credentials
3. Run the application:
```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8080`

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication Endpoints

#### Register
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Get Current User
```http
GET /auth/me
Authorization: Bearer <token>
```

#### Logout
```http
POST /auth/logout
Authorization: Bearer <token>
```

### Product Endpoints

#### Get All Products
```http
GET /products?category=electronics
```

#### Get Product by ID
```http
GET /products/:id
```

#### Create Product (Admin only)
```http
POST /products
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Product Name",
  "description": "Product description",
  "price": 99.99,
  "stock": 100,
  "category": "electronics",
  "image_url": "https://example.com/image.jpg"
}
```

#### Update Product (Admin only)
```http
PUT /products/:id
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Updated Name",
  "price": 89.99,
  "stock": 50
}
```

#### Delete Product (Admin only)
```http
DELETE /products/:id
Authorization: Bearer <admin-token>
```

### Cart Endpoints

#### Get Cart
```http
GET /cart
Authorization: Bearer <token>
```

#### Add Item to Cart
```http
POST /cart/items
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_id": 1,
  "quantity": 2
}
```

#### Update Cart Item
```http
PUT /cart/items/:itemId
Authorization: Bearer <token>
Content-Type: application/json

{
  "quantity": 3
}
```

#### Remove Item from Cart
```http
DELETE /cart/items/:itemId
Authorization: Bearer <token>
```

#### Clear Cart
```http
DELETE /cart
Authorization: Bearer <token>
```

### Order Endpoints

#### Create Order
```http
POST /orders
Authorization: Bearer <token>
Content-Type: application/json

{
  "shipping_address": "123 Main St, City, Country",
  "payment_method": "credit_card"
}
```

#### Get All Orders
```http
GET /orders
Authorization: Bearer <token>
```

#### Get Order by ID
```http
GET /orders/:id
Authorization: Bearer <token>
```

#### Update Order Status (Admin only)
```http
PUT /orders/:id/status
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "status": "shipped"
}
```

#### Cancel Order
```http
POST /orders/:id/cancel
Authorization: Bearer <token>
```

### User Endpoints

#### Get All Users
```http
GET /users
Authorization: Bearer <token>
```

#### Get User by ID
```http
GET /users/:id
Authorization: Bearer <token>
```

#### Update User
```http
PUT /users/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane@example.com"
}
```

#### Delete User (Admin only)
```http
DELETE /users/:id
Authorization: Bearer <admin-token>
```

### Post Endpoints

#### Get All Posts
```http
GET /posts
```

#### Get Post by ID
```http
GET /posts/:id
```

#### Create Post
```http
POST /posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Post Title",
  "content": "Post content here..."
}
```

#### Update Post
```http
PUT /posts/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Title",
  "content": "Updated content"
}
```

#### Delete Post
```http
DELETE /posts/:id
Authorization: Bearer <token>
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Server port | `8080` |
| `ENV` | Environment (development/staging/production) | `development` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `ecommerce` |
| `DB_SSLMODE` | Database SSL mode | `disable` |
| `JWT_SECRET` | JWT secret key | `your-secret-key-change-this` |
| `JWT_EXPIRE_TIME` | JWT expiration time (hours) | `24` |

## Database Schema

The application automatically creates the following tables:
- `users` - User accounts
- `products` - Product catalog
- `posts` - User posts
- `carts` - Shopping carts
- `cart_items` - Cart line items
- `orders` - Customer orders
- `order_items` - Order line items

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go
```

## Docker Commands

### Build and run development environment
```bash
docker-compose -f docker-compose.dev.yml up --build
```

### View logs
```bash
docker-compose -f docker-compose.dev.yml logs -f app
```

### Stop containers
```bash
docker-compose -f docker-compose.dev.yml down
```

### Remove volumes (fresh start)
```bash
docker-compose -f docker-compose.dev.yml down -v
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Support

For issues and questions, please create an issue in the repository.
