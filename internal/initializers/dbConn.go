package initializers

import (
	"log"
	"whotterre/argent/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB(connString string) {
	var err error
	DB, err = gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	if err := DB.AutoMigrate(&models.APIKey{}, &models.Transaction{}, &models.User{}, &models.Wallet{}); err != nil {
		log.Fatal("Failed to migrate database")
	}
	log.Println("Connected successfully to PostgreSQL database")
}

func GetDB() *gorm.DB {
	return DB
 }