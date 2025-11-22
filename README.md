# Chirpy 🐦

A production-ready Twitter-like social media API built with Go, featuring clean architecture, comprehensive testing, and modern development practices.

[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## ✨ Features

- 🔐 **User Authentication** - JWT-based auth with refresh tokens
- 📝 **Chirp Management** - Create, read, and delete chirps (140 char limit)
- 👤 **User Profiles** - User registration and profile updates
- 💳 **Premium Tier** - Chirpy Red subscription via webhooks
- 🚫 **Content Filtering** - Automatic profanity filtering
- 📊 **Admin Dashboard** - Metrics and system management
- 🧪 **Comprehensive Tests** - 38+ test cases, 80%+ coverage
- 🏗️ **Clean Architecture** - Layered design with clear separation
- 🐳 **Docker Support** - Easy local development
- 🔄 **CI/CD Ready** - GitHub Actions automated testing

## 🚀 Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL 14+

### Installation

```bash
# Clone the repository
git clone https://github.com/skylarhoughtongithub/chirpy.git
cd chirpy

# Install dependencies
go mod download

# Copy environment configuration
cp .env.example .env
# Edit .env with your configuration

# Run the application
make run
```

The server will start at `http://localhost:8080`

## 🔌 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/users` | Create user account | No |
| POST | `/api/login` | Login and get tokens | No |
| POST | `/api/refresh` | Refresh access token | Yes (Refresh) |
| POST | `/api/revoke` | Revoke refresh token | Yes (Refresh) |
| PUT | `/api/users` | Update user profile | Yes (Access) |

### Chirps

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/chirps` | Create a new chirp | Yes (Access) |
| GET | `/api/chirps` | Get all chirps | No |
| GET | `/api/chirps?author_id={id}` | Get chirps by author | No |
| GET | `/api/chirps?sort=desc` | Get chirps (newest first) | No |
| GET | `/api/chirps/{id}` | Get specific chirp | No |
| DELETE | `/api/chirps/{id}` | Delete chirp (owner only) | Yes (Access) |

### Webhooks

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/polka/webhooks` | Handle payment webhook | Yes (API Key) |

### Admin

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/admin/metrics` | View site metrics | No |
| POST | `/admin/reset` | Reset database (dev only) | No |

### Example API Usage

**Create a chirp:**
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"body":"Hello, world!"}'
```

**Get all chirps:**
```bash
curl http://localhost:8080/api/chirps
```

**Delete a chirp:**
```bash
curl -X DELETE http://localhost:8080/api/chirps/{chirp_id} \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run with coverage report
make test-coverage

# Run only unit tests (fast, no database)
go test ./tests/unit/...

# Run only integration tests
go test ./tests/integration/...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestChirpService_CreateChirp ./tests/unit/service/
```

## 🛠️ Development

### Available Commands

```bash
make help              # Show all commands
make build             # Build the application
make run               # Run the application
make test              # Run tests
make test-coverage     # Run tests with coverage
make lint              # Run linter
make format            # Format code
make clean             # Clean build artifacts
make sqlc              # Generate SQLC code
make migrate-up        # Run migrations up
make migrate-down      # Run migrations down
make migrate-status    # Check migration status
make docker-up         # Start Docker containers
make docker-down       # Stop Docker containers
make install-tools     # Install development tools
```

### Database Migrations

```bash
# Create new migration
goose -dir sql/schema create migration_name sql

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Check status
make migrate-status
```

### Adding New Features

1. **Define domain models** in `internal/domain/`
2. **Write SQL queries** in `sql/queries/`
3. **Generate code:** `make sqlc`
4. **Create repository** in `internal/repository/`
5. **Implement service** in `internal/service/`
6. **Create handler** in `internal/api/handlers/`
7. **Register route** in `internal/api/routes/`
8. **Write tests** in `tests/`

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DB_URL` | PostgreSQL connection string | Yes | - |
| `JWT_SECRET` | Secret key for JWT signing | Yes | - |
| `POLKA_KEY` | API key for Polka webhooks | Yes | - |
| `PORT` | Server port | No | 8080 |
| `PLATFORM` | Environment (dev/prod) | No | prod |

### Example `.env`

```bash
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key
POLKA_KEY=your-polka-api-key
PORT=8080
PLATFORM=dev
```

## 🐛 Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL is running
pg_isready -h localhost -p 5432

# Test connection
psql "postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"

# Check .env configuration
cat .env | grep DB_URL
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change port in .env
PORT=8081
```

### SQLC Generation Fails

```bash
# Install/update SQLC
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Regenerate
make sqlc
```

### Tests Failing

```bash
# Ensure test database exists
createdb chirpy_test

# Run migrations on test database
export TEST_DB_URL="postgres://postgres:postgres@localhost:5432/chirpy_test?sslmode=disable"
goose -dir sql/schema postgres "$TEST_DB_URL" up

# Clean and retry
make clean
make sqlc
make test
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⭐ Show Your Support

If you find this project helpful, please give it a star! ⭐

---
