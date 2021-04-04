package models

import (
	"fmt"
	"log"
	"time"
)

type Transaction struct {
	DateOfTransaction time.Time `gorm:"uniqueIndex:idx_unique_transaction"`
	Category          string
	SubCategory       string
	RuleDetected      *string
	Payee             string `gorm:"uniqueIndex:idx_unique_transaction"`
	ReversePayee      string
	CurrencyCode      string
	Amount            float32 `gorm:"uniqueIndex:idx_unique_transaction"`
	Out               bool    `gorm:"uniqueIndex:idx_unique_transaction"`
	Approved          bool
	Model
}

func (t Transaction) GetTransactionAsMap() map[string]string {
	ruleDetected := ""
	if t.RuleDetected == nil {
		ruleDetected = "None"
	} else {
		ruleDetected = *t.RuleDetected
	}
	return map[string]string{
		"Date of Transaction": t.DateOfTransaction.Format("January 2, 2006"),
		"Payee":               t.Payee,
		"Amount":              fmt.Sprintf("%.2f", t.Amount),
		"Approved":            "",
		"Category":            fmt.Sprintf("%s-%s", t.Category, t.SubCategory),
		"Rule Detected":       ruleDetected,
	}
}

func (m *DBManager) InsertTransactions(transactions []Transaction) error {
	result := m.Create(&transactions)
	log.Println(result)
	return result.Error
}
