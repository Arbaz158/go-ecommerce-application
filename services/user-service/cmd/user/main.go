package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-ecommerce-application/libs/kafka/config"
	"github.com/go-ecommerce-application/libs/kafka/consumer"
	profiling "github.com/go-ecommerce-application/libs/observability"
	"github.com/go-ecommerce-application/services/user-service/internal/database"
	"github.com/go-ecommerce-application/services/user-service/internal/handler"
	"github.com/go-ecommerce-application/services/user-service/internal/repository"
	"github.com/go-ecommerce-application/services/user-service/internal/routes"
	"github.com/go-ecommerce-application/services/user-service/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("../../../../.env")
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	addr := os.Getenv("HTTP_ADDR_USER_SERVICE")
	if addr == "" {
		addr = ":7071"
	}

	// set Gin mode via env; default to release for production
	if m := os.Getenv("GIN_MODE"); m != "" {
		gin.SetMode(m)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	db := database.ConnectMySQL()

	// build dependencies
	repo := repository.NewUserProfileRepository(db)
	svc := service.NewUserProfileService(repo)
	httpHandler := handler.NewUserProfileHandler(svc)
	kafkaEventHandler := handler.NewKafkaEventHandler(svc)

	// router and middleware
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// register routes
	routes.RegisterUserProfileRoutes(router, httpHandler)

	profiling.Start(profiling.Config{
		Enabled: os.Getenv("ENABLE_PPROF") == "true",
		Addr:    ":6061",
	})

	// Initialize Kafka consumer
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(kafkaBrokers) == 0 || kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}
	kafkaCfg := config.NewKafkaConfig(kafkaBrokers, "user-service-group")

	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	var consumerWg sync.WaitGroup

	// Start Kafka consumer in a goroutine
	kafkaConsumer, err := consumer.NewConsumer(kafkaCfg, "user.events", kafkaEventHandler.HandleUserSignedUpEvent)
	if err != nil {
		log.Printf("failed to initialize kafka consumer: %v (continuing without kafka)", err)
	} else {
		consumerWg.Add(1)
		go func() {
			defer consumerWg.Done()
			if err := kafkaConsumer.Start(consumerCtx); err != nil {
				log.Printf("kafka consumer error: %v", err)
			}
		}()
		log.Println("kafka consumer started")
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Printf("starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	// Stop Kafka consumer
	if kafkaConsumer != nil {
		consumerCancel()
		consumerWg.Wait()
		kafkaConsumer.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exiting")
}
