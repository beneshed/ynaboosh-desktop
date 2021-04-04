package models

type Category struct {
	YNABGroupID    string `gorm:"unique"`
	GroupName      string
	GroupHidden    bool
	GroupDeleted   bool
	YNABCategoryID string `gorm:"unique"`
	Name           string
	Hidden         bool
	Deleted        bool
	Model
}

func (db *DBManager) CreateCategories(categories []Category) {
	db.Create(&categories)
}
