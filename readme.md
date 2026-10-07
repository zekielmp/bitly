# E-Commerce API

A RESTful backend for an e-commerce platform, written in Go. It handles user authentication, product catalog management with image uploads, shopping carts, and orders.

## Features

- User registration and login with JWT access tokens and refresh tokens
- Product and category management
- Product image uploads (local disk or S3-compatible storage)
- Shopping cart and cart items
- Order creation and order history
- Event publishing via Watermill
- Dockerized local environment with Nginx and LocalStack

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Auth:** JWT
- **Events:** Watermill
- **Storage:** AWS S3 (LocalStack for local development) or local filesystem
- **Infrastructure:** Docker, Docker Compose, Nginx

## Project Structure

```
.
├── cmd/api/            # Application entry point
├── db/migrations/      # SQL migrations (up/down)
├── docker/             # Docker Compose, Nginx config, LocalStack init
├── internal/
│   ├── config/         # Configuration loading
│   ├── database/       # Database connection
│   ├── dto/            # Request/response data transfer objects
│   ├── interfaces/     # Abstractions (events, uploads)
│   ├── logger/         # Logging setup
│   ├── models/         # Domain models
│   ├── providers/      # Storage providers (local, S3/AWS)
│   ├── server/         # HTTP server, handlers, middleware
│   ├── services/       # Business logic
│   └── utils/          # JWT, password hashing, response helpers
├── uploads/            # Locally stored uploaded files
└── Makefile
```

## Getting Started

### Prerequisites

- Go 1.22+ (check `go.mod` for the exact version)
- Docker and Docker Compose
- Make

### Installation

```bash
git clone https://github.com/your-username/your-repo.git
cd your-repo
go mod download
```

### Configuration

Create a `.env` file in the project root:

```env
# Example, adjust to match internal/config/config.go
PORT= 8080
GIN_MODE=debug

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD= password
DB_NAME= bitly 
DB_SSLmODE = disable

JWT_SECRET = secret
JWT_EXPIRES_IN=24h
REFRESH_TOKEN_EXPIRES_IN=72h

AWS_REGION=east-1
AWS_ACCESS_KEY_ID =test
AWS_SECRET_ACCESS_KEY = test
AWS_S3_BUCKET= bitly-uploads
AWS_S3_ENDPOINT= https://Localstack:9000
UPLOAD_PROVIDER=local   # s3 or cld

UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE= 10485760 #10MB
```

### Start dependencies

```bash
# 1. Clone and install dependencies
git clone https://github.com/zekielmp/bitly.git
cd bitly
go mod download

# 2. Configure environment
cp .env.example .env   # then edit the values

# 3. Start PostgreSQL, LocalStack, and Nginx
make docker-up

# 4. Apply database migrations
make migrate-up

# 5. Run the API
make run
```

### Make commands

| Command             | Description                                |
|---------------------|--------------------------------------------|
| `make run`          | Run the API                                |
| `make build`        | Build the binary to `bin/app`              |
| `make lint`         | Format code and run golangci-lint          |
| `make format`       | Format code with `gofmt`                   |
| `make migrate-up`   | Apply database migrations                  |
| `make migrate-down` | Roll back database migrations              |
| `make docker-up`    | Start Docker services                      |
| `make docker-down`  | Stop Docker services                       |

## API Overview

## API Reference

Base URL: `http://localhost:<port>/api/v1`

**Auth levels:** 🌐 public, 🔒 requires a valid access token, 🛡️ requires admin role.

### Health & static files

| Method | Endpoint    | Auth | Description                         |
|--------|-------------|------|-------------------------------------|
| GET    | `/health`   | 🌐   | Health check (not under `/api/v1`)  |
| GET    | `/uploads/*`| 🌐   | Serves uploaded files (not under `/api/v1`) |

### Authentication

| Method | Endpoint         | Auth | Description               |
|--------|------------------|------|---------------------------|
| POST   | `/auth/register` | 🌐   | Create a new account      |
| POST   | `/auth/login`    | 🌐   | Log in and receive tokens |
| POST   | `/auth/refresh`  | 🌐   | Get a new access token    |
| POST   | `/auth/logout`   | 🌐   | Invalidate refresh token  |

### Users

| Method | Endpoint          | Auth | Description              |
|--------|-------------------|------|--------------------------|
| GET    | `/users/profile`  | 🔒   | Get current user profile |
| PUT    | `/users/profile`  | 🔒   | Update current profile   |

### Public catalog

| Method | Endpoint                | Auth | Description          |
|--------|-------------------------|------|----------------------|
| GET    | `/public/categories`    | 🌐   | List categories      |
| GET    | `/public/products`      | 🌐   | List products        |
| GET    | `/public/products/:id`  | 🌐   | Get a single product |

### Categories (admin)

| Method | Endpoint           | Auth | Description     |
|--------|--------------------|------|-----------------|
| POST   | `/categories`      | 🛡️   | Create category |
| PUT    | `/categories/:id`  | 🛡️   | Update category |
| DELETE | `/categories/:id`  | 🛡️   | Delete category |

### Products (admin)

| Method | Endpoint                | Auth | Description           |
|--------|-------------------------|------|-----------------------|
| POST   | `/products`             | 🛡️   | Create product        |
| PUT    | `/products/:id`         | 🛡️   | Update product        |
| DELETE | `/products/:id`         | 🛡️   | Delete product        |
| POST   | `/products/:id/images`  | 🛡️   | Upload product images |

### Cart

| Method | Endpoint           | Auth | Description          |
|--------|--------------------|------|----------------------|
| GET    | `/cart`            | 🔒   | Get current cart     |
| POST   | `/cart/items`      | 🔒   | Add item to cart     |
| PUT    | `/cart/items/:id`  | 🔒   | Update item quantity |
| DELETE | `/cart/items/:id`  | 🔒   | Remove item          |

### Orders

| Method | Endpoint       | Auth | Description              |
|--------|----------------|------|--------------------------|
| POST   | `/orders`      | 🔒   | Create order from cart   |
| GET    | `/orders`      | 🔒   | List your orders         |
| GET    | `/orders/:id`  | 🔒   | Get a single order       |

### Example

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret"}'
```

Send the returned access token on protected routes:

```bash
curl http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer <access_token>"
```

## Running Tests

```bash
go test ./...
```

## License

Distributed under the MIT License. See `LICENSE` for details.