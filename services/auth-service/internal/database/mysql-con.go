package database

import (
	"fmt"
	"log"
	"os"

	"github.com/go-ecommerce-application/services/auth-service/internal/domain/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func ConnectMySQL() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	authDBName := os.Getenv("AUTH_DB_NAME")

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, authDBName)
	DB, err = gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Println("Error while getting connection :", err)
		return
	}

	log.Println("Connection made to the database...")
	DB.AutoMigrate(&models.AuthUser{}, &models.RefreshToken{})
}
