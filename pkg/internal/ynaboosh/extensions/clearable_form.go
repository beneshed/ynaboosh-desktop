package extensions

import (
	"fyne.io/fyne/v2/widget"
)

type ClearableForm struct {
	*widget.Form
}

func NewClearableForm(form *widget.Form) *ClearableForm {
	f := &ClearableForm{
		form,
	}
	f.ExtendBaseWidget(f)
	return f
}

func (f *ClearableForm) Clear() {
	for _, item := range f.Items {
		typedItem, ok := item.Widget.(*widget.Entry)
		if ok {
			typedItem.SetText("")
		}
	}
	f.Refresh()
}
