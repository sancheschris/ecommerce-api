# Ecommerce API

A comprehensive RESTful API for managing an ecommerce platform built with Go. This API provides endpoints for managing products, orders, payments, and user authentication with integrated Stripe payment processing.

## Features

- **Product Management** - Create, read, update, and delete products
- **User Management** - User registration, authentication, and profile management
- **Order Management** - Complete order lifecycle management
- **Payment Processing** - Integrated Stripe payment processing
- **JWT Authentication** - Secure API access with JSON Web Tokens
- **API Documentation** - Interactive Swagger documentation
- **Comprehensive Testing** - Full test coverage for all components

## Tech Stack

- **Language**: Go 1.24+
- **Framework**: Chi HTTP router
- **Database**: GORM with SQLite (configurable)
- **Payment**: Stripe API integration
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI
- **Testing**: Go testing with testify

## Installation

### Prerequisites

- Go 1.24 or higher
- Git

### Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/ecommerce-api.git
   cd ecommerce-api
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   Create a `.env` file in the root directory:

   ```env
   # Server Configuration
   PORT=8080

   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=ecommerce_db

   # JWT Configuration
   JWT_SECRET=your-super-secret-jwt-key

   # Stripe Configuration
   STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
   STRIPE_PUBLISHABLE_KEY=pk_test_your_stripe_publishable_key
   ```

4. **Run database migrations**

   ```bash
   go run main.go migrate
   ```

5. **Start the server**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8080`

## Complete API Workflow

Follow this step-by-step guide to use the complete ecommerce flow:

### 1. Create a User Account

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 2. Generate Authentication Token

```bash
curl -X POST http://localhost:8080/users/generate_token \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

**Save the returned `access_token` for subsequent requests.**

### 3. Get or Create Products

List existing products:
```bash
curl -X GET http://localhost:8080/products \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Or create a new product:
```bash
curl -X POST http://localhost:8080/products \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Premium Laptop",
    "description": "High-performance laptop",
    "price": 1299.99,
    "stock_quantity": 50,
    "active": true
  }'
```

### 4. Create an Order

```bash
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [
      {
        "product_id": 1,
        "quantity": 1,
        "price": 1299.99
      }
    ]
  }'
```

### 5. Process Payment

```bash
curl -X POST http://localhost:8080/payments \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": 1,
    "amount_cents": 129999,
    "currency": "usd"
  }'
```

### 6. Check Payment Status

```bash
curl -X GET http://localhost:8080/payments/1/status \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Payment Statuses:**
- `requires_payment_method` - Waiting for payment method
- `processing` - Payment being processed  
- `succeeded` - Payment completed
- `canceled` - Payment canceled

## Testing

### Run all tests

```bash
go test ./...
```

### Run tests with coverage

```bash
go test -cover ./...
```

### Run tests for specific package

```bash
# Test payment repository
go test ./internal/repository/payment

# Test order models
go test ./internal/model

# Test handlers
go test ./internal/handler
```

### Run tests with verbose output

```bash
go test -v ./...
```

## API Documentation

Interactive API documentation is available via Swagger UI:

1. Start the server
2. Visit `http://localhost:8080/swagger/index.html`

The Swagger documentation provides:

- Complete endpoint reference
- Request/response schemas
- Interactive testing interface
- Authentication examples

## Project Structure

```
ecommerce-api/
├── cmd/
│   └── main.go
├── internal/
│   ├── handler/          # HTTP handlers
│   ├── model/           # Data models
│   ├── repository/      # Database layer
│   ├── service/         # Business logic
│   └── middleware/      # HTTP middleware
├── pkg/
│   └── stripe/          # External service clients
├── docs/               # Swagger documentation
├── migrations/         # Database migrations
└── tests/             # Integration tests
```

## Configuration

The application supports configuration via environment variables:

| Variable                 | Description            | Default     | Required |
| ------------------------ | ---------------------- | ----------- | -------- |
| `PORT`                   | Server port            | `8080`      | No       |
| `DB_HOST`                | Database host          | `localhost` | Yes      |
| `DB_PORT`                | Database port          | `5432`      | Yes      |
| `DB_USER`                | Database username      | -           | Yes      |
| `DB_PASSWORD`            | Database password      | -           | Yes      |
| `DB_NAME`                | Database name          | -           | Yes      |
| `JWT_SECRET`             | JWT signing secret     | -           | Yes      |
| `STRIPE_SECRET_KEY`      | Stripe secret key      | -           | Yes      |
| `STRIPE_PUBLISHABLE_KEY` | Stripe publishable key | -           | Yes      |

