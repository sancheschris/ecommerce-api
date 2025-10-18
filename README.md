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
- **Framework**: Gin HTTP framework
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

## Authentication

### Register a new user

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login to get access token

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### Using the token in requests

Include the token in the Authorization header:

```bash
curl -X GET http://localhost:8080/products \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## API Endpoints

### Products

- `GET /products` - List all products
- `GET /products/{id}` - Get product by ID
- `POST /products` - Create new product
- `PUT /products/{id}` - Update product
- `DELETE /products/{id}` - Delete product

### Orders

- `GET /orders` - List user orders
- `GET /orders/{id}` - Get order by ID
- `POST /orders` - Create new order

### Payments

- `POST /payments` - Process payment
- `GET /payments/{id}` - Get payment details

### Example: Create a product

```bash
curl -X POST http://localhost:8080/products \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "price": 999.99,
    "active": true
  }'
```

### Example: Create an order

```bash
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": 1,
        "quantity": 2,
        "unit_price": 999.99
      }
    ],
    "currency": "USD"
  }'
```

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

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Ensure all tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request
