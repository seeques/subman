# Subscription Service

REST API for aggregating user subscription data. Built with Go.

## Features

- CRUDL operations for subscriptions
- Calculate total subscription cost for a given period
- Filter by user ID and service name
- Pagination support
- Swagger documentation

## Tech Stack

- **Language:** Go 1.25
- **Router:** Chi
- **Database:** PostgreSQL 16
- **Migrations:** golang-migrate
- **Documentation:** Swagger (swaggo)
- **Containerization:** Docker

## Project Structure

```
.
├── main.go                 # Application entry point
├── internal/
│   ├── api/                # HTTP server setup
│   ├── config/             # Configuration loading
│   ├── handler/            # HTTP handlers
│   ├── models/             # Data models
│   ├── response/           # Response helpers
│   └── storage/            # Database operations
├── migrations/             # SQL migrations
├── docs/                   # Generated Swagger docs
├── Dockerfile
└── docker-compose.yml
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker and Docker Compose

### Configuration

1. Copy the example environment file:

```bash
cp .env.example .env
```

2. (Optional) Modify `.env` with your own values.

### Run with Docker (Recommended)

```bash
# Start all services
docker compose up -d --build

# View logs
docker compose logs -f app

# Stop services
docker compose down
```

### Run Locally

1. Copy environment file:

```bash
cp .env.example .env
```

2. Update `.env` for local development:

```env
DATABASE_URL=postgres://subscription:subscription@localhost:5436/subscription?sslmode=disable
```

3. Start PostgreSQL:

```bash
docker compose up -d postgres
```

4. Run migrations:

```bash
migrate -path ./migrations -database "postgres://subscription:subscription@localhost:5436/subscription?sslmode=disable" up
```

5. Run the application:

```bash
go run main.go
```

## API Documentation

Swagger UI available at: http://localhost:8080/swagger/index.html

## API Endpoints

| Method | Endpoint                           | Description            |
| ------ | ---------------------------------- | ---------------------- |
| POST   | `/api/v1/subscriptions`            | Create subscription    |
| GET    | `/api/v1/subscriptions`            | List all subscriptions |
| GET    | `/api/v1/subscriptions/{id}`       | Get subscription by ID |
| PUT    | `/api/v1/subscriptions/{id}`       | Update subscription    |
| DELETE | `/api/v1/subscriptions/{id}`       | Delete subscription    |
| GET    | `/api/v1/subscriptions/total-cost` | Calculate total cost   |

## Example Requests

### Create Subscription

```bash
curl -X POST "http://localhost:8080/api/v1/subscriptions" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

### List Subscriptions

```bash
curl "http://localhost:8080/api/v1/subscriptions?page=1&limit=10"
```

### Get Subscription

```bash
curl "http://localhost:8080/api/v1/subscriptions/1"
```

### Update Subscription

```bash
curl -X PUT "http://localhost:8080/api/v1/subscriptions/1" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 500,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025",
    "end_date": "12-2025"
  }'
```

### Delete Subscription

```bash
curl -X DELETE "http://localhost:8080/api/v1/subscriptions/1"
```

### Calculate Total Cost

```bash
curl "http://localhost:8080/api/v1/subscriptions/total-cost?start_period=01-2025&end_period=06-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
```

## Configuration

| Variable       | Description                  | Default |
| -------------- | ---------------------------- | ------- |
| `DATABASE_URL` | PostgreSQL connection string | -       |
| `PORT`         | Server port                  | 8080    |

## License

MIT
