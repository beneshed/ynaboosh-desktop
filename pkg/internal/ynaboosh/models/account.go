package models

type Account struct {
	YNABID  string `gorm:"unique"`
	Name    string
	Type    string
	Closed  bool
	Deleted bool
	Model
}

func (db *DBManager) CreateAccounts(accounts []Account) {
	db.Create(&accounts)
}
