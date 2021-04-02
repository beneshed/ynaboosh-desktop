package utils

import "fyne.io/fyne/v2/widget"

type Rule struct {
	Name  string
	Value string
}

type RuleList struct {
	Rules []Rule
	*widget.List
}

func NewRulesTable() {
	rules := make([]Rule, 0)
}
