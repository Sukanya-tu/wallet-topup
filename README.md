# Wallet Top-up System

A RESTful API built with Golang that supports verifying and confirming wallet top-up transactions. The application uses PostgreSQL for persistent storage, Redis for caching, JWT for authentication, and Docker for containerization.

---

## Run with Docker

### 1. Build and start all services

```bash
docker-compose up --build
```

### 2. Check if containers are running

```bash
docker ps
```

You should see these services:
- `wallet-topup-app` on port `8080`
- `wallet-topup-db` (PostgreSQL)
- `wallet-topup-redis` (Redis)

---

## Login to Get JWT Token

```http
POST /login
```

**Response:**

```json
{
  "token": "xxxxx.yyyyy.zzzzz"
}
```

Use the token as a Bearer token in `Authorization` header for all secured endpoints.

---

## API Endpoints

### Verify Top-up

```http
POST /api/verify
Authorization: Bearer <token>
```

**Request:**

```json
{
  "user_id": 1,
  "amount": 100.50,
  "payment_method": "credit_card"
}
```

**Response:**

```json
{
  "transaction_id": "abc123",
  "user_id": 1,
  "amount": 100.50,
  "payment_method": "credit_card",
  "status": "verified",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

---

### Confirm Top-up

```http
POST /api/confirm
Authorization: Bearer <token>
```

**Request:**

```json
{
  "transaction_id": "abc123"
}
```

**Response:**

```json
{
  "transaction_id": "abc123",
  "user_id": 1,
  "amount": 100.50,
  "status": "completed",
  "balance": 500.75
}
```

---

## Environment Variables

ใช้ `.env` ไฟล์ หรือใน `docker-compose.yml`:

```env
APP_ENV=dev
APP_PORT=8080
DB_HOST=db
DB_PORT=5432
DB_NAME=wallet-topup
DB_USER=postgres
DB_PASSWORD=password
DB_SSLMODE=disable
REDIS_ADDR=redis:6379
JWT_SECRET=myjwtsecretkey
USE_REAL_DB=true
```

---

## Features

- Golang + GORM (PostgreSQL)
- Redis Caching
- JWT Authentication
- RESTful API: `/verify` + `/confirm`
- Unit Testing (mock-based)
- Logging with logrus
- Docker + Docker Compose
- `.env` configuration support

---

## Run Unit Tests

```bash
go test ./... -v
```

> ใช้ mock ทั้ง Redis และ Database – ไม่ต้องเชื่อมต่อของจริง

---

## 🗂 Project Structure

```
wallet-topup/
- Dockerfile
- docker-compose.yml
- go.mod / go.sum
- .env
- main.go
- config/             # Env, DB, Redis, JWT
- handler/            # API handlers
- logs/               # Logger
- middleware/         # JWT auth middleware
- mocks/              # Mock interfaces for testing
- model/              # Structs + interfaces
- repository/         # GORM implementation
- service/            # Business logic
- wallet-log/         # Log output
```

---

## Notes

- ต้องสร้าง `users` ล่วงหน้าใน PostgreSQL (เช่น user_id=1)
- สามารถใช้ `docker exec -it wallet-topup-db psql -U postgres` เพื่อ insert test user
- ใช้งานได้ผ่าน Postman

---

## Author

Developed by Sukanya Tuamjun 
Assignment: Wallet Top-up System using Go + Docker + PostgreSQL + Redis + JWT
