package tables

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
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

type internalRulesTable struct {
	headers []string
	data    map[int]*models.Rule
}

func (r *internalRulesTable) Update(i widget.TableCellID, o fyne.CanvasObject) {
	if i.Row == 0 {
		if i.Col == 0 {
			o.(*widget.Label).SetText("Name")
		} else {
			o.(*widget.Label).SetText("Description")
		}
	} else {
		if i.Col == 0 {
			o.(*widget.Label).SetText(r.data[i.Row-1].Name)
		} else {
			o.(*widget.Label).SetText(r.data[i.Row-1].Description)
		}
	}

}

func (r *internalRulesTable) Length() (int, int) {
	return len(r.headers), len(r.data)
}

func NewRulesTable(form *extensions.ClearableForm, manager *models.DBManager) *RuleTable {
	initialState := manager.Rules()
	initialMap := make(map[int]*models.Rule)
	for i, rule := range initialState {
		initialMap[i] = &rule
	}
	internal := internalRulesTable{
		headers: []string{"Name", "Description"},
		data:    initialMap,
	}
	rulesTable := &RuleTable{
		internal,
		widget.NewTable(internal.Length, func() fyne.CanvasObject { return widget.NewLabel("") }, internal.Update),
	}
	rulesTable.ExtendBaseWidget(rulesTable)
	rulesTable.OnSelected = func(i widget.TableCellID) {
		log.Println("selected")
		//ruleEntry.SetText(state.Rules[i].Value)
		form.SubmitText = "Update"
		form.Refresh()
	}
	return rulesTable
}
