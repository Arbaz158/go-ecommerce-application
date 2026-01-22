package database

import (
	"log"

	"github.com/go-ecommerce-application/services/auth-service/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func ConnectMySQL() {
	url := "root:mysql123@tcp(localhost:3306)/auth_service"
	DB, err = gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Println("Error while getting connection :", err)
		return
	}

	log.Println("Connection made to the database...")
	DB.AutoMigrate(&models.AuthUser{}, &models.RefreshToken{})
}
