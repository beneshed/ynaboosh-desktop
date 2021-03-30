package institutions

import (
	"strconv"
	"strings"
	"time"
)

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

func parseSlashedDate(date string, seperator *string) time.Time {
	var dateParsed []string
	if seperator == nil {
		dateParsed = strings.Split(date, "/")
	} else {
		dateParsed = strings.Split(date, *seperator)
	}
	day, _ := strconv.Atoi(strings.TrimSpace(dateParsed[0]))
	month, _ := strconv.Atoi(strings.TrimSpace(dateParsed[1]))
	year, _ := strconv.Atoi(strings.TrimSpace(dateParsed[2]))
	return time.Date(year+2000, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
