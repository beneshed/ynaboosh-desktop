package forms

import (
	"bytes"
	"log"
	"text/template"

	"fyne.io/fyne/v2/widget"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/extensions"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/models"
)

const (
	exampleWhen            = `Transaction.Payee.Contains("לסרפוש")`
	exampleThen            = `example: Transaction.Category = "Groceries";Transaction.Payee = "Shufersal";`
	exampleRuleName        = "For example just put: Shufersal Groceries"
	exampleRuleDescription = "example: Assign grocery category to shufersal and switch to english"
	genericRuleTemplate    = `rule {{.Name}} "{{.Description}}" salience {{.Priority}} {
		when
			{{.When}}
		then
			{{.Then}}
		}
	`
)

var ruleTemplate = template.Must(template.New("rule").Parse(genericRuleTemplate))

type RawRule struct {
	Name        string
	Description string
	When        string
	Then        string
	Priority    int64
}

func (r RawRule) ToRule() *models.Rule {
	var rule bytes.Buffer
	_ = ruleTemplate.Execute(&rule, r)
	return &models.Rule{
		Name:        r.Name,
		Description: r.Description,
		When:        r.When,
		Then:        r.Then,
		Priority:    10,
	}
}

func NewRuleEditForm(manager *models.DBManager) *extensions.ClearableForm {
	whenEntry := widget.NewMultiLineEntry()
	whenEntryFormItem := widget.NewFormItem("When", whenEntry)
	whenEntryFormItem.HintText = exampleWhen
	thenEntry := widget.NewMultiLineEntry()
	thenFormEntry := widget.NewFormItem("Then", thenEntry)
	thenFormEntry.HintText = exampleThen
	ruleNameEntry := widget.NewEntry()
	ruleNameFormItem := widget.NewFormItem("Rule Name", ruleNameEntry)
	ruleNameFormItem.HintText = exampleRuleName
	descriptionEntry := widget.NewEntry()
	descriptionFormItem := widget.NewFormItem("Description", descriptionEntry)
	descriptionFormItem.HintText = exampleRuleDescription

	createRuleForm := widget.NewForm(
		ruleNameFormItem,
		descriptionFormItem,
		whenEntryFormItem,
		thenFormEntry,
	)
	clearableCreateRuleForm := extensions.NewClearableForm(createRuleForm)
	clearableCreateRuleForm.OnSubmit = func() {
		rawRule := RawRule{
			Name:        ruleNameEntry.Text,
			Description: descriptionEntry.Text,
			When:        whenEntry.Text,
			Then:        thenEntry.Text,
			Priority:    10,
		}

		rule := rawRule.ToRule()
		_, err := manager.FindOrCreateRuleFromRule(rule)
		log.Println(err)
		clearableCreateRuleForm.Clear()
	}
	clearableCreateRuleForm.SubmitText = "Add"
	clearableCreateRuleForm.OnCancel = func() {
		clearableCreateRuleForm.SubmitText = "Add"
		clearableCreateRuleForm.Clear()

	}
	return clearableCreateRuleForm
}
