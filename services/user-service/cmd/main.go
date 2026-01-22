package main

import (
	"log"

	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-ecommerce-application/services/internal/profiling"
	"github.com/go-ecommerce-application/services/user-service/internal/database"
	"github.com/go-ecommerce-application/services/user-service/internal/handler"
	"github.com/go-ecommerce-application/services/user-service/internal/repository"
	"github.com/go-ecommerce-application/services/user-service/internal/routes"
	"github.com/go-ecommerce-application/services/user-service/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8081"
	}

	// set Gin mode via env; default to release for production
	if m := os.Getenv("GIN_MODE"); m != "" {
		gin.SetMode(m)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	db := database.ConnectMySQL()

	// // Load JWT configuration
	// jwtConfig, err := auth.LoadFromEnv()
	// if err != nil {
	// 	log.Fatalf("Failed to load JWT config: %v", err)
	// }

	// build dependencies
	// tokenManager := auth.NewTokenManager(jwtConfig)
	repo := repository.NewUserProfileRepository(db)
	svc := service.NewUserProfileService(repo)
	h := handler.NewUserProfileHandler(svc)

	// router and middleware
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// register routes with JWT middleware
	routes.RegisterUserProfileRoutes(router, h)

	profiling.Start(profiling.Config{
		Enabled: os.Getenv("ENABLE_PPROF") == "true",
		Addr:    ":6061",
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
