# 💰 Wallet Top-up API (Golang)

A RESTful API for verifying and confirming wallet top-up transactions, built using **Golang**, **GORM**, **Redis**, and **Docker** as part of an assignment project.

---

## 📌 Features

- ✅ Verify top-up transaction (`POST /wallet/verify`)
- ✅ Confirm top-up transaction (`POST /wallet/confirm`)
- ✅ GORM for database interaction
- ✅ Redis for caching verified transactions
- ✅ Logging with logrus
- ✅ Dockerized for easy deployment
- ✅ Unit tests with coverage

---

## 🚀 Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/sukanya_tu/wallet-topup.git
cd wallet-topup
```

---

## 🐳 Run with Docker

> Requires Docker & Docker Compose installed

```bash
docker-compose up --build
```

The API will be available at: [http://localhost:8080](http://localhost:8080)

---

## 🔧 Environment Variables

These are set in `docker-compose.yml`, but you can override them:


## 🔌 API Endpoints

### `POST /wallet/verify`

**Request Body:**
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

### `POST /wallet/confirm`

**Request Body:**
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

## 🧪 Run Unit Tests

```bash
# If SQLite is used for testing
CGO_ENABLED=1 go test ./... -v

# Optional: with coverage
go test ./... -cover
```

> If you’re on Windows and using SQLite, make sure you have GCC installed (e.g. via TDM-GCC)

---

## 🧩 Notes

- Project uses SQLite for fast unit testing and PostgreSQL + Redis in Docker.
- Tests cover service and handler layers using GORM, Redis mock, and logrus.
- Project structure follows basic Hexagonal Architecture.

---

## 📬 Author

- Sukanya Tuamjun
- Assignment Submission: Wallet Top-up System (Golang)
