package db

import (
	"gorm.io/gorm"
	"time"
)

type HealthCounter struct {
	CheckId  uint      `gorm:"primaryKey;autoIncrement"`
	Datetime time.Time `gorm:"not null"`
}

type Database interface {
	Create(value interface{}) *gorm.DB
}
