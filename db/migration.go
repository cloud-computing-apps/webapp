package db

import (
	"gorm.io/gorm"
	"log"
)

func InitDB(dbConnection *gorm.DB) {
	err := dbConnection.AutoMigrate(&HealthCounter{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	log.Println("Database migrated successfully!")
}
