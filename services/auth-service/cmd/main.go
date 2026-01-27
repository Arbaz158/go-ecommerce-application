package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/go-ecommerce-application/pkg/kafka/config"
	"github.com/go-ecommerce-application/pkg/kafka/producer"
	"github.com/go-ecommerce-application/services/auth-service/internal/database"
	"github.com/go-ecommerce-application/services/auth-service/internal/handler"
	"github.com/go-ecommerce-application/services/auth-service/internal/repository"
	"github.com/go-ecommerce-application/services/auth-service/internal/routes"
	"github.com/go-ecommerce-application/services/auth-service/internal/service"
	"github.com/go-ecommerce-application/services/internal/profiling"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":7070"
	}

	// set Gin mode via env; default to release for production
	if m := os.Getenv("GIN_MODE"); m != "" {
		gin.SetMode(m)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	database.ConnectMySQL()

	// Initialize Kafka producer
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(kafkaBrokers) == 0 || kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}
	kafkaCfg := config.NewKafkaConfig(kafkaBrokers, "")
	kafkaProducer, err := producer.NewProducer(kafkaCfg)
	if err != nil {
		log.Printf("failed to initialize kafka producer: %v (continuing without kafka)", err)
	}
	defer func() {
		if kafkaProducer != nil {
			kafkaProducer.Close()
		}
	}()

	// build dependencies
	repo := repository.NewAuthRepository()
	svc := service.NewAuthService(repo, kafkaProducer)
	h := handler.NewAuthHandler(svc)

	// router and middleware
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// register routes
	routes.RegisterAuthRoutes(router, h)

	profiling.Start(profiling.Config{
		Enabled: os.Getenv("ENABLE_PPROF") == "true",
		Addr:    ":6060",
	})

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exiting")
}
