package models

import (
	"time"
	"fmt"
)

type Transaction struct {
	DateOfTransaction time.Time `sql:"unique_index:idx_date_payee_amount"`
	Payee             string `sql:"unique_index:idx_date_payee_amount"`
	CurrencyCode      string
	Amount            float32 `sql:"unique_index:idx_date_payee_amoutn"`
	Out               bool
	Approved          bool
	Model
}

func (t Transaction) GetTransactionAsMap() map[string]string {
	return map[string]string{
		"Date of Transaction": t.DateOfTransaction.Format("January 2, 2006"),
		"Payee":               t.Payee,
		"Amount":              fmt.Sprintf("%.2f", t.Amount),
		"Approved": "",
	}
}

func NewTransation() *Transaction {
	return &Transaction{}
}