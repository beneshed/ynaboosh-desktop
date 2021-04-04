package utils

import 
(
	"fyne.io/fyne/v2/widget"
)

type indexedForm struct {
	Index map[int]*widget.FormItem
	*widget.Form
}

func newIndexedForm(index map[int]*widget.FormItem) *indexedForm {
	form := &indexedForm{
		index,
		widget.NewForm(),
	}
	for _, item := range index {
		form.AppendItem(item)
	}
	form.ExtendBaseWidget(form)
	return form
}