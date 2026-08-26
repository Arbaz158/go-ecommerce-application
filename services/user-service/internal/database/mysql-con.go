package database

import (
	"fmt"
	"log"

	"os"

	"github.com/go-ecommerce-application/services/user-service/internal/domain/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQL() *gorm.DB {

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	userDBName := os.Getenv("USER_DB_NAME")

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, userDBName)
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Println("error while getting connection :", err)
		return nil
	}
	log.Println("Database connected successfully")

	db.AutoMigrate(&models.UserProfile{}, &models.Address{})
	return db

}
