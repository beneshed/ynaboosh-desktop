package extensions

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Rule struct {
	Name  string
	Value string
}

type RuleListState struct {
	Rules []Rule
}

func (r RuleListState) Length() int {
	return len(r.Rules)
}

func (r *RuleListState) Update(i widget.ListItemID, o fyne.CanvasObject) {
	o.(*widget.Label).SetText(r.Rules[i].Name)
}

type RuleList struct {
	RuleListState
	*widget.List
}

func NewRulesList(form *ClearableForm) *RuleList {
	state := RuleListState{
		[]Rule{
			{
				Name:  "foo",
				Value: "bar",
			},
			{
				Name:  "foo",
				Value: "bar",
			},
			{
				Name:  "foo",
				Value: "bar",
			},
			{
				Name:  "foo",
				Value: "bar",
			},
			{
				Name:  "foo",
				Value: "bar",
			},
			{
				Name:  "foo",
				Value: "bar",
			},
			{
				Name:  "foo",
				Value: "bar",
			},
			{
				Name:  "foo",
				Value: "bar",
			},
			{
				Name:  "foo",
				Value: "bar",
			},
		},
	}
	rulesList := &RuleList{
		state,
		widget.NewList(state.Length, func() fyne.CanvasObject {
			return widget.NewLabel("template")
		}, state.Update),
	}
	rulesList.ExtendBaseWidget(rulesList)
	rulesList.OnSelected = func(i widget.ListItemID) {
		log.Println("selected")
		//ruleEntry.SetText(state.Rules[i].Value)
		form.SubmitText = "Update"
		form.Refresh()
	}
	return rulesList
}
