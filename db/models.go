package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type HealthCounter struct {
	CheckId  uint      `gorm:"primaryKey;autoIncrement"`
	Datetime time.Time `gorm:"not null"`
}

type FileTable struct {
	Id         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FileName   string    `gorm:"not null" json:"file_name"`
	Url        string    `gorm:"not null" json:"url"`
	UploadDate time.Time `gorm:"not null;autoCreateTime" json:"upload_date"`
}

type Database interface {
	Create(value any) error
	Delete(value any, conds ...any) error
	FindByID(id uuid.UUID, out any) error
}

type GormDatabase struct {
	DB *gorm.DB
}

func (g *GormDatabase) Create(value any) error {
	return g.DB.Create(value).Error
}

func (g *GormDatabase) Delete(value any, conds ...any) error {
	return g.DB.Delete(value, conds...).Error
}

func (g *GormDatabase) FindByID(id uuid.UUID, out any) error {
	return g.DB.First(out, "id = ?", id).Error
}
