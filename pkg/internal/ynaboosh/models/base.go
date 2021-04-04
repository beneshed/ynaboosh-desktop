package models

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID uint `gorm:"primaryKey"`
	CreatedOn time.Time
	UpdatedOn time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}