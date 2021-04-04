package main

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

type Setting struct {
	AccessToken *string
	ExpiresOn   time.Time
	Model
}

type Account struct {
	YNABID  string
	Name    string
	Type    string
	Closed  bool
	Deleted bool
	Model
}

type Category struct {
	YNABGroupID    string
	GroupName      string
	GroupHidden    bool
	GroupDeleted   bool
	YNABCategoryID string
	Name           string
	Hidden         bool
	Deleted        bool
	Model
}

func InitializeDB(db *gorm.DB) error {
	err := db.AutoMigrate(&Setting{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Account{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Category{})
	if err != nil {
		return err
	}

	return nil
}