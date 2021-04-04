package ynaboosh

import (
	"gorm.io/gorm"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/models"
)

func InitializeDB(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Account{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Category{})
	if err != nil {
		return err
	}

	return nil
}