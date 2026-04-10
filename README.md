# 🐦 Chirpy — Go HTTP Server

A JSON REST API server built in Go with a PostgreSQL backend. Created while working through the [HTTP Servers in Go](https://www.boot.dev/courses/learn-http-servers-golang) course from [Boot.dev](https://boot.dev).

Chirpy is a simple social platform API where users can register, authenticate, and post short messages called "chirps" (think: tweets, but limited to 140 characters). It covers core backend concepts like JWT-based authentication, refresh token management, database migrations, webhook handling, and clean JSON API design.

For full API documentation, see the [docs folder](./docs).

---

## Prerequisites

- **Go** (1.22+) — https://go.dev/doc/install
- **PostgreSQL** (v15+) — https://www.postgresql.org/download/
- **Goose** (database migration tool):
```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
```

---

## Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/AggroSec/Go-HTTP-Server.git
cd Go-HTTP-Server
```

### 2. Set up the database

Connect to PostgreSQL and create the database:

```bash
# Linux / WSL
sudo -u postgres psql

# macOS (Homebrew)
psql postgres
```

```sql
CREATE DATABASE chirpy;
\q
```

### 3. Run the migrations

```bash
goose postgres "postgres://postgres:[password]@localhost:5432/chirpy?sslmode=disable" -dir sql/schema up
```

Replace `[password]` with your PostgreSQL password.

### 4. Configure environment variables

The server reads its configuration from environment variables. Create a `.env` file in the project root:

```env
DB_URL="postgres://postgres:[password]@localhost:5432/chirpy?sslmode=disable"
JWT_SECRET=your_secret_key_here
POLKA_KEY=your_webhook_api_key_here
PLATFORM="dev"
```

- `JWT_SECRET` can be any long random string — it's used to sign and verify access tokens.
- `POLKA_KEY` is the API key used to authenticate incoming webhook events from the Polka payment service.
- `PLATFORM` should be set to `dev` for local development; this enables certain admin endpoints.

### 5. Build and run

```bash
go build -o chirpy && ./chirpy
```

The server will start on port `8080` by default. You can verify it's running by hitting the health endpoint:

```bash
curl http://localhost:8080/api/healthz
```

---

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **DB Migrations:** [Goose](https://github.com/pressly/goose)
- **SQL Code Generation:** [sqlc](https://sqlc.dev/)
- **Auth:** JWT access tokens + refresh tokens (argon2id password hashing)

---

## Documentation

Full API documentation including all endpoints, request/response formats, and authentication details can be found in the [docs folder](./docs).

---

## Acknowledgements

Built following the [HTTP Servers in Go](https://www.boot.dev/courses/build-web-server-golang) course from [Boot.dev](https://boot.dev) — a great hands-on platform for learning backend development.

---

*Licensed under MIT*
