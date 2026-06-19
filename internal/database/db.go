package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/turpatur/dcg-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() *gorm.DB {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	db.AutoMigrate(&models.DBCard{}, &models.PairBan{})

	return db
}
