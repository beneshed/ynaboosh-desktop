package ynabimporter

import _ "embed"

type FinancialInstitution string

var (
	//go:embed assets/discount_he.json
	EmbeddedDiscountHE string
)

const (
	Discount FinancialInstitution = "Discount"
)
