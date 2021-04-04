package forms

import (
	"bytes"
	"log"
	"text/template"

	"fyne.io/fyne/v2/widget"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/extensions"
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
	Priority    int
}

func NewRuleEditForm() *extensions.ClearableForm {
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
		whenEntryFormItem,
		descriptionFormItem,
		ruleNameFormItem,
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
		var rule bytes.Buffer
		err := ruleTemplate.Execute(&rule, rawRule)
		if err != nil {
			log.Panicln(err)
		}
		log.Println(rule.String())
		clearableCreateRuleForm.Clear()
	}
	clearableCreateRuleForm.SubmitText = "Add"
	clearableCreateRuleForm.OnCancel = func() {
		clearableCreateRuleForm.SubmitText = "Add"
		clearableCreateRuleForm.Clear()

	}
	return clearableCreateRuleForm
}
