package tables

import (
	"encoding/json"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	grule "github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/extensions"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/models"
)

type Rule struct {
	Name  string
	Value string
}

type RuleTable struct {
	internalRulesTable
	*widget.Table
}
type ImportRulesAction struct {
	FileSelected *string
	Dialog       *dialog.FileDialog
}

func BuildRawRules(data []byte) []grule.GruleJSON {
	var rules []grule.GruleJSON
	_ = json.Unmarshal(data, &rules)
	return rules
}

func (r *RuleTable) ImportRules() {
	window := fyne.CurrentApp().Driver().AllWindows()[0]
	var fileSelected *string
	action := ImportRulesAction{
		FileSelected: fileSelected,
		Dialog: dialog.NewFileOpen(func(info fyne.URIReadCloser, err error) {
			log.Println(info, err)
			if err == nil && info != nil {
				pathPtr := info.URI().Path()
				fileSelected = &pathPtr
				bytes, err := os.ReadFile(pathPtr)
				log.Println(string(bytes), err)
				knowledgeLibrary := ast.NewKnowledgeLibrary()
				ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
				ruleset, err := grule.ParseJSONRuleset(bytes)
				rules := BuildRawRules(bytes)
				r.db.InsertRawRules(rules)
				log.Println(ruleset, err)
				parsedRules := grule.NewBytesResource([]byte(ruleset))
				err = ruleBuilder.BuildRuleFromResource("ImportedRules", "0.0.1", parsedRules)
				log.Println(err)
			}
		}, window),
	}
	action.Dialog.Show()
	r.Table.Refresh()
}

func (r *RuleTable) WrapTableWidth() {
	var maxNameWidth float32 = fyne.MeasureText(r.headers[0], theme.TextSize(), fyne.TextStyle{}).Width
	var maxDescriptionWidth float32 = fyne.MeasureText(r.headers[1], theme.TextSize(), fyne.TextStyle{}).Width
	r.Table.SetColumnWidth(0, maxNameWidth+roundPad)
	r.Table.SetColumnWidth(1, maxDescriptionWidth+roundPad)
	r.Table.Refresh()
}

type internalRulesTable struct {
	headers []string
	data    map[int]*models.Rule
	db      models.DBManager
}

func (r *internalRulesTable) Update(i widget.TableCellID, o fyne.CanvasObject) {
	log.Println(i.Row, i.Col, r.data)
	if i.Row == 0 {
		if i.Col == 0 {
			o.(*widget.Label).SetText("Name")
		} else {
			o.(*widget.Label).SetText("Description")
		}
	} else {
		if i.Col == 0 {
			o.(*widget.Label).SetText(r.data[i.Row].Name)
		} else {
			o.(*widget.Label).SetText(r.data[i.Row].Description)
		}
	}

}

/*
Length Rows by Columns
*/
func (r *internalRulesTable) Length() (int, int) {
	return len(r.data) + 1, len(r.headers)
}

func NewRulesTable(form *extensions.ClearableForm, manager *models.DBManager) *RuleTable {
	initialState := manager.Rules()
	log.Println(initialState)
	initialMap := make(map[int]*models.Rule)
	for i, rule := range initialState {
		initialMap[i] = &rule
	}
	internal := internalRulesTable{
		headers: []string{"Name", "Description"},
		data:    initialMap,
		db:      *manager,
	}
	rulesTable := &RuleTable{
		internal,
		widget.NewTable(internal.Length, func() fyne.CanvasObject { return widget.NewLabel("") }, internal.Update),
	}
	rulesTable.WrapTableWidth()
	rulesTable.ExtendBaseWidget(rulesTable)
	rulesTable.OnSelected = func(i widget.TableCellID) {
		log.Println("selected")
		//ruleEntry.SetText(state.Rules[i].Value)
		form.SubmitText = "Update"
		form.Refresh()
	}
	return rulesTable
}
