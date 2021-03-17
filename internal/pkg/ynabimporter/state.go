package ynabimporter

import (
	"go.bmvs.io/ynab"
	"gorm.io/gorm"
)

type GlobalState struct {
	DB         *gorm.DB
	YNABClient ynab.ClientServicer
}

func NewGlobalState() *GlobalState {
	return &GlobalState{}
}
