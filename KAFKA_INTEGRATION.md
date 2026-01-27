# Kafka Integration Guide

## Directory Structure

```
pkg/kafka/
├── config/
│   └── kafka.go           # Kafka configuration and broker validation
├── events/
│   └── events.go          # Event definitions (UserSignedUp, etc.)
├── producer/
│   └── producer.go        # Generic Kafka producer
└── consumer/
    └── consumer.go        # Generic Kafka consumer
```

## Architecture

This is a production-ready, scalable Kafka integration that follows microservices architecture principles:

- **Centralized Kafka Package**: All Kafka functionality lives in `pkg/kafka/` making it reusable across all services
- **Event Definitions**: All events are defined in `pkg/kafka/events/` for consistency
- **Service Isolation**: Each service independently initializes producer/consumer based on its role
- **Error Handling**: Graceful handling of Kafka unavailability (services continue without Kafka)

## Current Flow

### 1. Auth Service (Producer)

**Location**: `services/auth-service/`

- When a user signs up, the `Signup` method publishes a `UserSignedUp` event to the `user.events` topic
- Event contains: UserID, Email, FirstName, LastName, EventID, Timestamp

**Configuration**:
```env
KAFKA_BROKERS=localhost:9092
```

### 2. User Service (Consumer)

**Location**: `services/user-service/`

- Consumer subscribes to `user.events` topic with consumer group `user-service-group`
- Upon receiving `UserSignedUp` event, automatically creates a user profile entry
- Prevents duplicate profiles via UserID uniqueness check

**Configuration**:
```env
KAFKA_BROKERS=localhost:9092
```

## Usage Examples

### Adding a New Producer Service

```go
import (
    "github.com/go-ecommerce-application/pkg/kafka/config"
    "github.com/go-ecommerce-application/pkg/kafka/producer"
)

// In main.go
kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
kafkaCfg := config.NewKafkaConfig(kafkaBrokers, "")
kafkaProducer, err := producer.NewProducer(kafkaCfg)

// In service
kafkaProducer.Publish(ctx, "topic.name", key, messageBytes)
```

### Adding a New Consumer Service

```go
import (
    "github.com/go-ecommerce-application/pkg/kafka/config"
    "github.com/go-ecommerce-application/pkg/kafka/consumer"
)

// Define handler function
handler := func(ctx context.Context, message []byte) error {
    // Process message
    return nil
}

// In main.go
kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
kafkaCfg := config.NewKafkaConfig(kafkaBrokers, "service-consumer-group")
kafkaConsumer, err := consumer.NewConsumer(kafkaCfg, "topic.name", handler)

// Start consuming (blocks)
kafkaConsumer.Start(ctx)
```

### Adding New Events

1. Define in `pkg/kafka/events/events.go`:
```go
type OrderPlaced struct {
    EventType string
    EventID   string
    OrderID   string
    UserID    string
    Amount    float64
    Timestamp time.Time
}

func (e *OrderPlaced) ToJSON() ([]byte, error) {
    return json.Marshal(e)
}
```

2. Publish from producer service
3. Consume in consumer service

## Starting Local Kafka

```bash
# Using Docker Compose (create docker-compose.yml)
version: '3'
services:
  kafka:
    image: confluentinc/cp-kafka:7.5.0
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    depends_on:
      - zookeeper
  
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

# Run
docker-compose up -d
```

## Testing

All services handle Kafka unavailability gracefully. If Kafka is not available:
- Services continue to operate normally
- Events are simply not published/consumed
- No errors are thrown at service startup

## Future Enhancements

- Dead Letter Queue (DLQ) for failed message processing
- Event schema registry for versioning
- Message retry logic with exponential backoff
- Event ordering guarantees
- Monitoring and metrics collection
