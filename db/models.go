package db

import (
	"time"
)

type HealthCounter struct {
	CheckId  uint      `gorm:"primaryKey;autoIncrement"`
	Datetime time.Time `gorm:"not null"`
}
