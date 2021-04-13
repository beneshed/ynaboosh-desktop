package ynaboosh

import (
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/models"
	"gorm.io/gorm"
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
	err = db.AutoMigrate(&models.Rule{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Transaction{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.RawRule{})

	return nil
}
