package models

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