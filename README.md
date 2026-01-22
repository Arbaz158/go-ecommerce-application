# Go E-Commerce Application

A microservices-based e-commerce platform built with Go, featuring a modular architecture with shared authentication, comprehensive testing, and production-ready profiling capabilities.

## 📋 Table of Contents

- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Services](#services)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Profiling](#profiling)
- [Development Guide](#development-guide)
- [Contributing](#contributing)

---

## 🏗️ Project Overview

This is a modular microservices architecture for an e-commerce platform with:

- **Authentication Service** (`auth-service`): Handles user registration, login, token management
- **User Service** (`user-service`): Manages user profiles and addresses
- **Shared Authentication Package** (`pkg/auth`): Reusable JWT and middleware for all services
- **Comprehensive Testing**: Unit tests for repositories, services, and handlers
- **Production Profiling**: Built-in CPU, memory, and goroutine profiling with pprof

### Key Features

✅ JWT-based authentication with access & refresh tokens  
✅ Password hashing with bcrypt  
✅ Middleware-based route protection  
✅ Database abstraction with GORM  
✅ Mock-based unit testing  
✅ HTTP handler testing with Gin  
✅ SQL mocking for database tests  
✅ Production-ready pprof profiling  
✅ Graceful shutdown  

---

## 🏛️ Architecture

### Microservices Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     API Gateway / Client                    │
└────────────────────────┬────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
        ▼                ▼                ▼
   ┌──────────┐     ┌──────────┐    ┌─────────┐
   │   Auth   │     │   User   │    │  Other  │
   │ Service  │     │ Service  │    │Services │
   └────┬─────┘     └────┬─────┘    └─────────┘
        │                │
        └────────────────┼────────────────┐
                         │                │
                    ┌────▼─────┐    ┌─────▼──────┐
                    │ pkg/auth  │    │   MySQL    │
                    │(Shared)   │    │  Database  │
                    └───────────┘    └────────────┘
```

### Service Layer Pattern

Each service follows a clean architecture with:

```
Handler (HTTP Layer)
    ↓
Service (Business Logic)
    ↓
Repository (Data Layer)
    ↓
Database
```

### Shared Authentication

All services use the centralized `pkg/auth` package for:
- JWT token generation and validation
- Password hashing and verification
- Auth middleware for route protection
- Shared DTOs (Data Transfer Objects)

---

## 📦 Prerequisites

### Required

- **Go 1.25.1** or higher
- **MySQL 8.0** or higher
- **git**

### Optional

- **Postman** or **curl** for API testing
- **go-profiler-ui** for viewing profiling results

### Installation

```bash
# Install Go (macOS)
brew install go

# Install MySQL (macOS)
brew install mysql
```

---

## 📁 Project Structure

```
go-ecommerce-application/
├── pkg/
│   └── auth/                          # Shared authentication package
│       ├── jwt.go                     # JWT token generation & validation
│       ├── utils.go                   # Password hashing utilities
│       ├── dto.go                     # Shared data transfer objects
│       └── middleware.go              # Gin authentication middleware
│
├── services/
│   ├── auth-service/
│   │   ├── cmd/
│   │   │   └── main.go                # Service entry point
│   │   └── internal/
│   │       ├── handler/               # HTTP handlers
│   │       ├── service/               # Business logic
│   │       ├── repository/            # Data access layer
│   │       ├── models/                # Database models
│   │       ├── dto/                   # Local DTOs (deprecated - use pkg/auth)
│   │       ├── authentication/        # JWT logic (deprecated - use pkg/auth)
│   │       ├── database/              # Database configuration
│   │       ├── routes/                # Route definitions
│   │       └── utils/                 # Utilities (deprecated - use pkg/auth)
│   │
│   ├── user-service/
│   │   ├── cmd/
│   │   │   └── main.go                # Service entry point
│   │   └── internal/
│   │       ├── handler/               # HTTP handlers
│   │       │   └── user-profile-handlers_test.go
│   │       ├── service/               # Business logic
│   │       │   └── user-profile-service_test.go
│   │       ├── repository/            # Data access layer
│   │       │   └── user-profile-repository_test.go
│   │       ├── models/                # Database models
│   │       ├── database/              # Database configuration
│   │       ├── routes/                # Route definitions
│   │       └── utils/                 # Utilities
│   │
│   └── internal/
│       └── profiling/                 # Shared profiling configuration
│           └── pprof.go               # pprof initialization
│
├── go.mod                             # Go modules file
├── go.sum                             # Module checksums
├── PROFILING_GUIDE.md                 # Profiling documentation
├── PROFILING_CHEATSHEET.md            # Quick profiling reference
├── test-profiling.sh                  # Script to generate profiles
└── README.md                          # This file
```

---

## 🚀 Services

### 1. Auth Service

**Purpose**: User authentication and authorization

**Port**: `8080` (configurable via `HTTP_ADDR`)

**Endpoints**:
- `POST /auth/signup` - Register a new user
- `POST /auth/login` - Login user (returns access & refresh tokens)
- `POST /auth/refresh` - Refresh access token
- `GET /auth/logout` - Logout user (protected)

**Database Models**:
```go
type AuthUser struct {
    Id       string  // UUID
    Email    string  // Unique
    Password string  // Bcrypt hashed
    Role     string
    Status   string  // "active", "inactive", etc.
}

type RefreshToken struct {
    Id        string
    UserID    string
    ExpiresAt time.Time
}
```

### 2. User Service

**Purpose**: User profile and address management

**Port**: `8081` (configurable via `HTTP_ADDR`)

**Endpoints**:
- `GET /health` - Health check (no auth required)
- `GET /users/me` - Get user profile (protected)
- `POST /users/address` - Create address (protected)
- `GET /users/address` - Get user addresses (protected)

**Database Models**:
```go
type UserProfile struct {
    ID    uint
    Name  string
    Phone string  // Unique
    Email string  // Unique
}

type Address struct {
    ID         uint
    UserID     uint      // Foreign key
    Street     string
    City       string
    State      string
    PostalCode string
    CreatedAt  int64
    UpdatedAt  int64
}
```

---

## 🏁 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/your-org/go-ecommerce-application.git
cd go-ecommerce-application
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Configure Environment

Create a `.env` file in the root directory:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ecommerce_db

# Service Configuration
HTTP_ADDR=:8080          # Auth service
GIN_MODE=release         # or "debug" for development

# Profiling
ENABLE_PPROF=true        # Enable profiling on port 6060, 6061
```

### 4. Setup Database

```bash
# Start MySQL server
brew services start mysql

# Create database
mysql -u root -p < schema.sql

# Or run migrations manually
mysql -u root -p -D ecommerce_db < migrations/001_initial_schema.sql
```

### 5. Run Services

**Terminal 1 - Auth Service**:
```bash
cd services/auth-service/cmd
go run main.go
```

**Terminal 2 - User Service**:
```bash
cd services/user-service/cmd
go run main.go
```

Both services should start without errors. Check logs for confirmation.

---

## 📡 API Documentation

### Authentication Flow

#### 1. Register User

```bash
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword",
    "role": "user"
  }'
```

**Response** (201):
```json
{
  "message": "Signup successful"
}
```

#### 2. Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

**Response** (200):
```json
{
  "response": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "JWT",
    "expires_in": 15,
    "refresh_expires_in": 10080,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "role": "user",
      "status": "active"
    }
  }
}
```

#### 3. Access Protected Endpoint

```bash
curl -X GET http://localhost:8081/users/me \
  -H "Authorization: Bearer <access_token>"
```

**Response** (200):
```json
{
  "id": 1,
  "name": "John Doe",
  "phone": "123-456-7890",
  "email": "john@example.com"
}
```

#### 4. Create Address (Protected)

```bash
curl -X POST http://localhost:8081/users/address \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "postal_code": "10001"
  }'
```

**Response** (201):
```json
{
  "message": "Address created successfully"
}
```

---

## 🧪 Testing

### Run All Tests

```bash
# Run all tests in the project
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

### Run Tests by Package

```bash
# Auth Service Tests
go test -v github.com/go-ecommerce-application/services/auth-service/internal/repository
go test -v github.com/go-ecommerce-application/services/auth-service/internal/service

# User Service Tests
go test -v github.com/go-ecommerce-application/services/user-service/internal/repository
go test -v github.com/go-ecommerce-application/services/user-service/internal/service
go test -v github.com/go-ecommerce-application/services/user-service/internal/handler
```

### Run Specific Test

```bash
# Run a single test
go test -run TestGetUserProfile_Success -v github.com/go-ecommerce-application/services/user-service/internal/service
```

### Test Coverage Report

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View HTML coverage report
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

### Test Structure

Tests use **mock implementations** to isolate layers:

- **Repository Tests** (`*_test.go`): Mock database with `sqlmock`
- **Service Tests** (`*_test.go`): Mock repository interface
- **Handler Tests** (`*_test.go`): Mock service interface + Gin test context

Example service test:
```go
func TestGetUserProfile_Success(t *testing.T) {
    mockRepo := &MockUserProfileRepository{
        GetUserProfileByIDFunc: func(id uint) (*models.UserProfile, error) {
            return &models.UserProfile{ID: 1, Name: "John"}, nil
        },
    }
    
    service := NewUserProfileService(mockRepo)
    profile, err := service.GetUserProfile(1)
    
    if err != nil || profile.Name != "John" {
        t.Fatalf("Test failed")
    }
}
```

---

## 📊 Profiling

This application includes production-ready profiling using Go's `pprof`.

### Enable Profiling

Set environment variable:
```bash
export ENABLE_PPROF=true
```

### Available Profiles

- **CPU Profile**: `/debug/pprof/profile` (30-second capture)
- **Memory Profile**: `/debug/pprof/heap`
- **Goroutine Profile**: `/debug/pprof/goroutine`
- **Block Profile**: `/debug/pprof/block`
- **Mutex Profile**: `/debug/pprof/mutex`

### Capture CPU Profile

```bash
# Capture 30-second CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Interactive analysis
(pprof) top10     # Top 10 functions by CPU time
(pprof) list main # Show source code for main
(pprof) web       # Generate graph (requires graphviz)
```

### Memory Profiling

```bash
# Capture heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# View allocations
(pprof) alloc_space  # Total allocations
(pprof) inuse_space  # Current in-use memory
(pprof) top
```

### Goroutine Profiling

```bash
# Capture goroutine profile
go tool pprof http://localhost:6060/debug/pprof/goroutine
(pprof) top
```

See [PROFILING_GUIDE.md](PROFILING_GUIDE.md) for detailed profiling instructions.

---

## 👨‍💻 Development Guide

### Adding New Endpoint

1. **Define Model** (`internal/models/`):
```go
type NewEntity struct {
    ID   uint
    Name string
}
```

2. **Create Repository** (`internal/repository/`):
```go
type NewEntityRepository interface {
    Create(entity *NewEntity) error
    Get(id uint) (*NewEntity, error)
}
```

3. **Implement Service** (`internal/service/`):
```go
func (s *service) Create(entity NewEntity) error {
    return s.repo.Create(&entity)
}
```

4. **Create Handler** (`internal/handler/`):
```go
func (h *handler) CreateEntity(c *gin.Context) {
    // Handle HTTP request
}
```

5. **Register Route** (`internal/routes/`):
```go
func RegisterRoutes(r *gin.Engine, h *handler.Handler) {
    r.POST("/entities", h.CreateEntity)
}
```

6. **Write Tests** (`*_test.go`):
```go
func TestCreate_Success(t *testing.T) {
    // Test implementation
}
```

### Using Shared Auth Package

```go
import "github.com/go-ecommerce-application/pkg/auth"

// Generate tokens
accessToken, refreshToken, _, _, err := auth.GenerateTokens(userID, role)

// Validate token
claims, err := auth.ValidateAccessToken(tokenString)

// Hash password
hashedPassword, err := auth.HashPassword(password)

// Verify password
isValid := auth.CheckPasswordHash(password, hashedPassword)

// Protect routes
router.GET("/protected", auth.AuthMiddleware(), handler.ProtectedHandler)
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | Service port |
| `GIN_MODE` | `release` | `debug` or `release` |
| `DB_HOST` | `localhost` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL user |
| `DB_PASSWORD` | - | MySQL password |
| `DB_NAME` | `ecommerce_db` | Database name |
| `ENABLE_PPROF` | `false` | Enable profiling |

---

## 🔒 Security

### Best Practices Implemented

✅ **Password Security**
- Bcrypt hashing with default cost factor
- Never stored in plain text
- Salted hashes

✅ **JWT Tokens**
- HS256 signing algorithm
- 15-minute access token expiry
- 7-day refresh token expiry
- Claims include UserID and Role

✅ **Route Protection**
- Middleware-based authentication
- Bearer token validation
- Context-based user extraction

✅ **Input Validation**
- JSON schema validation
- Type checking
- Error handling

### Secrets Management

Store secrets in `.env` file (not version controlled):
```env
DB_PASSWORD=secure_password_here
JWT_SECRET=your_jwt_secret_here
```

---

## 🐛 Troubleshooting

### Database Connection Failed

```
Error: Error 1045: Access denied for user 'root'@'localhost'
```

**Solution**: Check `.env` file credentials match MySQL configuration.

### Port Already in Use

```
listen: bind: address already in use
```

**Solution**: Change port in `.env` or kill existing process:
```bash
lsof -i :8080
kill -9 <PID>
```

### JWT Token Invalid

```
Error: invalid access token
```

**Solution**: Ensure token is passed in `Authorization: Bearer <token>` header.

### Tests Failing

```bash
# Run tests with output
go test -v ./...

# Check for mock setup issues
# Verify mock function fields are initialized
```

---

## 📚 Key Dependencies

| Dependency | Version | Purpose |
|-----------|---------|---------|
| `gin-gonic/gin` | v1.11.0 | HTTP web framework |
| `gorm.io/gorm` | v1.31.1 | ORM for database |
| `gorm.io/driver/mysql` | v1.6.0 | MySQL driver |
| `golang-jwt/jwt` | v5 | JWT token generation |
| `golang.org/x/crypto` | v0.40.0 | Bcrypt password hashing |
| `google/uuid` | v1.6.0 | UUID generation |
| `joho/godotenv` | v1.5.1 | `.env` file parsing |
| `DATA-DOG/go-sqlmock` | v1.5.2 | SQL mocking for tests |

---

## 📝 License

This project is licensed under the MIT License.

---

## 🤝 Contributing

### Code Style

- Follow standard Go conventions
- Use `go fmt` for formatting
- Write descriptive variable names
- Add comments for exported functions

### Pull Request Process

1. Create feature branch: `git checkout -b feature/your-feature`
2. Write tests for new functionality
3. Ensure all tests pass: `go test ./...`
4. Commit changes: `git commit -am 'Add your message'`
5. Push to branch: `git push origin feature/your-feature`
6. Create Pull Request

---

## 📞 Support

For issues or questions:

1. Check [PROFILING_GUIDE.md](PROFILING_GUIDE.md) for profiling help
2. Review test examples in `*_test.go` files
3. Check logs for error messages
4. Open an issue on GitHub

---

## 🎯 Roadmap

- [ ] Add order service
- [ ] Implement payment integration
- [ ] Add product catalog service
- [ ] Implement API gateway
- [ ] Add Docker support
- [ ] Kubernetes deployment manifests
- [ ] OpenAPI/Swagger documentation
- [ ] GraphQL API layer

---

**Created**: January 2026  
**Last Updated**: January 22, 2026
