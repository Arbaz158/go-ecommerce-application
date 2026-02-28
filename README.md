# Go E-Commerce Application

A microservices-based e-commerce platform built with Go, featuring event-driven architecture with Kafka, centralized shared libraries, comprehensive testing, and production-ready profiling capabilities.

## 📋 Table of Contents

- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Services](#services)
- [Event-Driven Architecture](#event-driven-architecture)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Profiling](#profiling)
- [Performance Testing](#performance-testing)
- [Development Guide](#development-guide)
- [Contributing](#contributing)

---

## 🏗️ Project Overview

This is a modular microservices architecture for an e-commerce platform with:

- **Authentication Service** (`auth-service`): Handles user registration, login, token management, publishes user signup events
- **User Service** (`user-service`): Manages user profiles and addresses, consumes user signup events
- **Shared Libraries** (`libs/`): Centralized authentication, Kafka, and observability packages
  - `libs/auth`: JWT and middleware for all services
  - `libs/kafka`: Event producer/consumer for inter-service communication
  - `libs/observability`: Profiling and monitoring utilities
- **Comprehensive Testing**: Unit tests for repositories, services, and handlers
- **Production Profiling**: Built-in CPU, memory, and goroutine profiling with pprof
- **Event-Driven Architecture**: Kafka-based asynchronous communication between services

### Key Features

✅ JWT-based authentication with access & refresh tokens  
✅ Password hashing with bcrypt  
✅ Middleware-based route protection  
✅ Database abstraction with GORM  
✅ Event-driven inter-service communication with Kafka  
✅ Mock-based unit testing  
✅ HTTP handler testing with Gin  
✅ SQL mocking for database tests  
✅ Production-ready pprof profiling  
✅ Graceful shutdown and signal handling  

---

## 🏛️ Architecture

### Microservices Architecture with Event-Driven Flow

```
┌──────────────────────────────────────────────────────────────────┐
│                     API Gateway / Client                         │
└────────────────┬─────────────────────────┬──────────────────────┘
                 │                         │
        ┌────────▼──────────┐      ┌───────▼────────┐
        │  Auth Service     │      │  User Service  │
        │  (Port 7070)      │      │  (Port 7071)   │
        └────────┬──────────┘      └───────┬────────┘
                 │                         │
                 │    ┌──────────────┐     │
                 ├───▶│ Kafka Topic  ◀─────┤
                 │    │user.events   │     │
                 │    └──────────────┘     │
                 │                         │
        ┌────────▼──────────────────────────▼────────┐
        │         Shared Libraries (libs/)           │
        ├────────────────────────────────────────────┤
        │  • libs/auth       - JWT & Middleware      │
        │  • libs/kafka      - Producer/Consumer     │
        │  • libs/observability - Profiling & PPprof │
        └────────────────────────────────────────────┘
                 │
        ┌────────▼──────────┐
        │  MySQL Database   │
        │  (Shared)         │
        └───────────────────┘
```

### Service Layer Pattern

Each service follows a clean architecture with request flow:

```
HTTP Request
    ↓
Handler (HTTP Layer) - Route mapping & request validation
    ↓
Service (Business Logic) - Core domain logic
    ↓
Repository (Data Layer) - Database abstraction
    ↓
Database - MySQL persistence
```

### Event-Flow Communication

```
Auth Service               Kafka                User Service
    │                       │                        │
    ├─→ User Signup ───→ user.events ───→ Consume Event
    │    Published          Topic          Handle Event
    │                                      Create Profile
    │                                      (Auto sync)
```

### Shared Libraries

All services use centralized `libs/` packages:

- **`libs/auth`**: JWT token generation/validation, password hashing, auth middleware
- **`libs/kafka`**: Event producer/consumer, Kafka configuration, event definitions
- **`libs/observability`**: pprof profiling, performance monitoring

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
├── libs/                                  # Shared libraries (go modules)
│   ├── auth/                              # Authentication library
│   │   ├── go.mod                         # Auth module definition
│   │   ├── jwt.go                         # JWT token generation & validation
│   │   ├── utils.go                       # Password hashing utilities
│   │   ├── dto.go                         # Shared data transfer objects
│   │   └── middleware.go                  # Gin authentication middleware
│   │
│   ├── kafka/                             # Event-driven communication library
│   │   ├── go.mod                         # Kafka module definition
│   │   ├── config/
│   │   │   └── kafka.go                   # Kafka configuration & broker setup
│   │   ├── events/
│   │   │   └── events.go                  # Event definitions (UserSignedUp, etc.)
│   │   ├── producer/
│   │   │   └── producer.go                # Generic Kafka producer
│   │   └── consumer/
│   │       └── consumer.go                # Generic Kafka consumer with group handling
│   │
│   └── observability/                     # Observability library
│       ├── go.mod                         # Observability module definition
│       └── pprof.go                       # pprof profiling initialization
│
├── services/
│   ├── auth-service/                      # Authentication & Authorization service
│   │   ├── go.mod
│   │   ├── cmd/
│   │   │   └── auth/
│   │   │       └── main.go                # Service entry point (Port 7070)
│   │   ├── Dockerfile
│   │   ├── docs/
│   │   └── internal/
│   │       ├── handler/
│   │       │   └── auth-handlers.go       # HTTP request handlers
│   │       ├── service/
│   │       │   └── auth-service.go        # Business logic (signup, login, refresh)
│   │       ├── repository/
│   │       │   ├── auth-repo.go           # Database operations
│   │       │   └── auth-repo_test.go      # Repository unit tests
│   │       ├── models/
│   │       │   └── user.go                # Database models (User, RefreshToken)
│   │       ├── domain/
│   │       │   ├── authentication/        # Authentication domain logic
│   │       │   │   └── jwt.go
│   │       │   ├── dto/
│   │       │   │   └── auth-dto.go        # Request/response DTOs
│   │       │   └── utils/
│   │       │       └── auth-utils.go
│   │       ├── routes/
│   │       │   └── auth-routes.go         # Route registration
│   │       └── database/
│   │           └── mysql-con.go           # MySQL connection setup
│   │
│   └── user-service/                      # User Profile service
│       ├── go.mod
│       ├── cmd/
│       │   └── user/
│       │       └── main.go                # Service entry point (Port 7071)
│       ├── Dockerfile
│       ├── docs/
│       └── internal/
│           ├── handler/
│           │   ├── user-profile-handlers.go       # HTTP endpoints
│           │   ├── user-profile-handlers_test.go  # Handler tests
│           │   └── kafka-event-handler.go         # Kafka event consumers
│           ├── service/
│           │   ├── user-profile-service.go        # Business logic
│           │   └── user-profile-service_test.go   # Service tests
│           ├── repository/
│           │   ├── user-profile-repository.go     # Database operations
│           │   └── user-profile-repository_test.go # Repository tests
│           ├── models/
│           │   └── user_profile.go                # Database models (UserProfile, Address)
│           ├── domain/
│           │   └── dto/
│           ├── routes/
│           │   └── user-profile-routes.go         # Route registration
│           ├── utils/
│           └── database/
│               └── mysql-con.go           # MySQL connection setup
│
├── go.work                                # Go workspace file (for local module development)
├── go.work.sum
├── auth-tester.go                         # Load testing tool for auth service
├── PROFILING_GUIDE.md                     # Detailed profiling documentation
├── PROFILING_CHEATSHEET.md                # Quick profiling reference
├── KAFKA_INTEGRATION.md                   # Kafka event architecture documentation
├── KAFKA_BUILD_STATUS.md                  # Kafka setup verification guide
├── README.md                              # This file
└── docs/                                  # Additional documentation
```

### Key Structure Notes

- **Go Workspace**: Uses `go.work` to manage multiple modules (auth, kafka, observability are separate modules)
- **Internal Package**: Each service uses `internal/` to hide implementation details
- **Shared Libraries**: Located in `libs/` and imported as separate modules by services
- **Testing**: Tests placed alongside implementation files with `*_test.go` suffix
- **Database Models**: Separate from DTOs (models are DB-focused, DTOs are API-focused)

---

## 🚀 Services

### 1. Auth Service

**Purpose**: User authentication, authorization, and token management

**Port**: `7070` (configurable via `HTTP_ADDR`)

**Endpoints**:
- `POST /auth/signup` - Register a new user (publishes `UserSignedUp` event to Kafka)
- `POST /auth/login` - Login user (returns access & refresh tokens)
- `POST /auth/refresh` - Refresh access token
- `GET /auth/logout` - Logout user (protected)

**Responsibilities**:
- User registration and credential validation
- JWT token generation and management
- Kafka event publishing for user signup events
- Password hashing and verification

### 2. User Service

**Purpose**: User profile and address management, consumes auth events

**Port**: `7071` (configurable via `HTTP_ADDR_USER_SERVICE`)

**Endpoints**:
- `GET /health` - Health check (no auth required)
- `GET /users/me` - Get user profile (protected)
- `POST /users/address` - Create address (protected)
- `GET /users/address` - Get user addresses (protected)

**Responsibilities**:
- User profile management
- Address management for users
- Kafka event consumption (UserSignedUp events)
- Automatic user profile creation from signup events

---

## 📨 Event-Driven Architecture

This application uses **Kafka** for asynchronous, event-driven communication between services.

### Current Event Flow

**Scenario**: User Registration

```
1. User calls PUT /auth/signup on Auth Service
    ↓
2. Auth Service validates credentials and creates user
    ↓
3. Auth Service publishes UserSignedUp event to Kafka topic 'user.events'
    ↓
4. User Service consumes UserSignedUp event
    ↓
5. User Service automatically creates UserProfile entry
    ↓
6. Services remain decoupled - no direct HTTP calls needed
```

### Event Definitions

Located in `libs/kafka/events/events.go`:

```go
type UserSignedUp struct {
    EventType string    // "user.signed_up"
    UserID    string    // User ID from auth service
    Email     string
    FirstName string
    LastName  string
    EventID   string    // Unique event ID
    Timestamp int64     // Unix timestamp
}
```

### Using the Event System

**Publishing an Event** (Auth Service):
```go
// Publish user signup event to Kafka
event := events.UserSignedUp{
    EventType: "user.signed_up",
    UserID:    userID,
    Email:     user.Email,
    Timestamp: time.Now().Unix(),
}

messageBytes, _ := json.Marshal(event)
kafkaProducer.Publish(ctx, "user.events", userID, messageBytes)
```

**Consuming an Event** (User Service):
```go
// Consumer automatically calls this handler for each message
handler := func(ctx context.Context, message []byte) error {
    var event events.UserSignedUp
    json.Unmarshal(message, &event)
    
    // Create user profile from event
    return userProfileService.HandleUserSignedUpEvent(&event)
}

consumer, _ := consumer.NewConsumer(kafkaCfg, "user.events", handler)
consumer.Start(ctx) // Blocks and consumes messages
```

### Benefits of This Architecture

✅ **Loose Coupling**: Services don't call each other directly  
✅ **Scalability**: Easy to add new services that consume events  
✅ **Resilience**: Services continue working if Kafka is temporarily unavailable  
✅ **Asynchronous Processing**: Events don't block the API response  
✅ **Audit Trail**: All events are logged in Kafka for debugging  

For detailed Kafka setup and configuration, see [KAFKA_INTEGRATION.md](KAFKA_INTEGRATION.md).

---

## 🏁 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/your-org/go-ecommerce-application.git
cd go-ecommerce-application
```

### 2. Install Dependencies

```bash
# Download and verify all module dependencies
go mod download

# Tidy dependencies
go mod tidy
```

### 3. Setup Database

```bash
# Start MySQL server (macOS)
brew services start mysql

# Create database
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS ecommerce_db;"

# (Optional) Run migrations if schema.sql exists
# mysql -u root -p ecommerce_db < schema.sql
```

### 4. Setup Kafka (Optional but Recommended)

For event-driven features to work, Kafka must be running:

```bash
# Install Kafka using Homebrew
brew install kafka

# Start Kafka broker (runs on localhost:9092 by default)
brew services start kafka

# Verify Kafka is running
kafka-broker-api-versions --bootstrap-server localhost:9092
```

If Kafka is unavailable, services will log warnings but continue operating (graceful degradation).

### 5. Configure Environment

Create a `.env` file in the root directory:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_mysql_password
DB_NAME=ecommerce_db

# Service Configuration
HTTP_ADDR=:7070                    # Auth service port
HTTP_ADDR_USER_SERVICE=:7071       # User service port
GIN_MODE=release                   # or "debug" for development

# Kafka Configuration
KAFKA_BROKERS=localhost:9092       # Kafka broker addresses (comma-separated)

# Profiling (Optional)
ENABLE_PPROF=true                  # Enable pprof profiling on port 6060, 6061
```

### 6. Run Services

**Terminal 1 - Auth Service**:
```bash
cd services/auth-service/cmd/auth
go run main.go
# Expected output: Server listening on :7070
```

**Terminal 2 - User Service**:
```bash
cd services/user-service/cmd/user
go run main.go
# Expected output: Server listening on :7071
# Should also show: Listening for Kafka events...
```

Both services should start without errors. Check console logs for confirmation:
- Auth service should log: `listening to address :7070`
- User service should log: `listening to address :7071` and Kafka consumer status

---

## 📡 API Documentation

### Authentication Flow

#### 1. Register User

```bash
curl -X POST http://localhost:7070/auth/signup \
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

**Background**: This endpoint also publishes a `UserSignedUp` event to Kafka, which is consumed by the User Service to automatically create a user profile.

#### 2. Login

```bash
curl -X POST http://localhost:7070/auth/login \
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
curl -X GET http://localhost:7071/users/me \
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
curl -X POST http://localhost:7071/users/address \
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

## � Performance Testing

This project includes an integrated load testing tool to measure auth service performance.

### Using auth-tester.go

Located at the root: `auth-tester.go`

This tool simulates concurrent load on the Auth Service login endpoint:

**Configuration**:
```go
const (
    url         = "http://localhost:7070/auth/login"
    concurrency = 50    // number of goroutines
    requests    = 2000  // total requests
)
```

**Run performance test**:
```bash
go run auth-tester.go
```

**What it tests**:
- Sends 2000 concurrent login requests (50 goroutines)
- Measures total execution time
- Tests under realistic concurrent load
- Helps identify performance bottlenecks

**Example output**:
```
Completed 2000 requests in 12.34 seconds
Requests/sec: 162
Avg latency: ~6ms
```

**Customize the test**:
```bash
# Edit auth-tester.go to change:
# - concurrency level
# - total number of requests
# - timeout duration
# - endpoint URL
# - request payload
```

### Integration with Profiling

Combine performance testing with profiling for insights:

```bash
# Terminal 1: Start services with profiling enabled
export ENABLE_PPROF=true
cd services/auth-service/cmd/auth && go run main.go

# Terminal 2: Start performance test and capture CPU profile simultaneously
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=20 &
go run auth-tester.go

# Terminal 3: Analyze heap after test
go tool pprof http://localhost:6060/debug/pprof/heap
```

---

## �👨‍💻 Development Guide

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

### Using Shared Libraries

#### Using Shared Auth Library

```go
import "github.com/go-ecommerce-application/libs/auth"

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

#### Using Kafka Producer/Consumer

```go
import (
    "github.com/go-ecommerce-application/libs/kafka/config"
    "github.com/go-ecommerce-application/libs/kafka/producer"
    "github.com/go-ecommerce-application/libs/kafka/consumer"
)

// Publishing events
kafkaCfg := config.NewKafkaConfig(brokers, "")
kafkaProducer, _ := producer.NewProducer(kafkaCfg)
kafkaProducer.Publish(ctx, "topic.name", key, messageBytes)

// Consuming events
handler := func(ctx context.Context, message []byte) error {
    // Process message
    return nil
}
kafkaCfg := config.NewKafkaConfig(brokers, "consumer-group")
consumer, _ := consumer.NewConsumer(kafkaCfg, "topic.name", handler)
consumer.Start(ctx)
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:7070` | Auth service port |
| `HTTP_ADDR_USER_SERVICE` | `:7071` | User service port |
| `GIN_MODE` | `release` | `debug` or `release` |
| `DB_HOST` | `localhost` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL user |
| `DB_PASSWORD` | - | MySQL password |
| `DB_NAME` | `ecommerce_db` | Database name |
| `KAFKA_BROKERS` | `localhost:9092` | Comma-separated Kafka broker addresses |
| `ENABLE_PPROF` | `false` | Enable profiling (pprof on port 6060, 6061) |

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
# Check which process is using port 7070
lsof -i :7070

# Kill the process
kill -9 <PID>

# Or change the port in .env:
# HTTP_ADDR=:7072
```

### Kafka Connection Failed

```
Error: Failed to connect to Kafka brokers: localhost:9092
Error: context deadline exceeded
```

**Solution**: 
- Ensure Kafka is running: `brew services start kafka`
- Verify Kafka broker is accessible: `kafka-broker-api-versions --bootstrap-server localhost:9092`
- Check `KAFKA_BROKERS` environment variable
- Note: Services will continue operating without Kafka (graceful degradation)

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
| `segmentio/kafka-go` | v0.4.50 | Kafka client for events |

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
**Last Updated**: February 28, 2026
