package models

import (
	"encoding/json"
	"log"

	grule "github.com/hyperjumptech/grule-rule-engine/pkg"
)

type Rule struct {
	Name        string `gorm:"unique"`
	Description string
	Definition  string
	When        string
	Then        string
	Priority    int64 `gorm:"default:10"`
	Model
}

type RawRule struct {
	Name    string `gorm:"unique"`
	Body    string
	Version string
}

func (db *DBManager) FindOrCreateRule(name string, description string, definition string, when string, then string, priority *int64) (*Rule, error) {
	rule := &Rule{
		Name:        name,
		Description: description,
		Definition:  definition,
		When:        when,
		Then:        then,
		Priority:    10,
	}
	result := db.Where(Rule{Name: name}).FirstOrCreate(rule)
	if result.Error != nil {
		return nil, result.Error
	}
	return rule, nil
}

func (db *DBManager) FindOrCreateRuleFromRule(rule *Rule) (*Rule, error) {
	result := db.Where(Rule{Name: rule.Name}).FirstOrCreate(rule)
	if result.Error != nil {
		return nil, result.Error
	}
	return rule, nil
}

func (db *DBManager) FindRule(name string) (*Rule, error) {
	rule := &Rule{}
	result := db.Where("name = ?", name).First(rule)
	if result.Error != nil {
		return nil, result.Error
	}
	return rule, nil
}

func (db *DBManager) Rules() []Rule {
	var rules []Rule
	_ = db.Find(&rules)
	return rules
}

func (db *DBManager) InsertRawRules(rules []grule.GruleJSON) {
	var rawRules []RawRule
	for _, rule := range rules {
		rawRuleBytes, _ := json.Marshal(rule)
		rawRules = append(rawRules, RawRule{
			Name:    rule.Name,
			Body:    string(rawRuleBytes),
			Version: "0.0.1",
		})
	}
	result := db.Create(rawRules)
	log.Println(result)
}
