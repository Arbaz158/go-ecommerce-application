package database

import (
	"log"

	"github.com/go-ecommerce-application/services/user-service/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQL() *gorm.DB {

	url := "root:mysql123@tcp(localhost:3306)/user_service?parseTime=true"
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Println("error while getting connection :", err)
		return nil
	}
	log.Println("Database connected successfully")

	db.AutoMigrate(&models.UserProfile{}, &models.Address{})
	return db

}
