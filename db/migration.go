package db

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitDB(dbConnection *gorm.DB) {
	err := dbConnection.AutoMigrate(&HealthCounter{}, &FileTable{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	log.Info("Database migrated successfully!")
}
