package models

type Account struct {
	YNABID  string
	Name    string
	Type    string
	Closed  bool
	Deleted bool
	Model
}