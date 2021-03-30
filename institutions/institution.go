package institutions

import "time"

type Institution int

const (
	nis               = "₪"
	usd               = "$"
	CalHE Institution = iota
	DiscountEN
	DiscountHE
)

func (i Institution) String() string {
	return [...]string{"CalHE", "DiscountEN", "DiscountHE"}[i]
}

type Transaction struct {
	DateOfTransaction time.Time
	Payee             string
	CurrencyCode      string
	Amount            float32
	Out               bool
}

func NewTransation() *Transaction {
	return &Transaction{}
}
