package models

import (
	"fmt"
	"time"
)

type Transaction struct {
	DateOfTransaction time.Time `sql:"unique_index:idx_date_payee_amount"`
	Category          string
	SubCategory       string
	RuleDetected      *string
	Payee             string `sql:"unique_index:idx_date_payee_amount"`
	CurrencyCode      string
	Amount            float32 `sql:"unique_index:idx_date_payee_amoutn"`
	Out               bool
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

func NewTransation() *Transaction {
	return &Transaction{}
}
