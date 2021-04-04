package models

import "gorm.io/gorm"

type DBManager struct {
	*gorm.DB
}
