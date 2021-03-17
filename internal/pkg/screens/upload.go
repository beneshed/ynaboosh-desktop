package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewUploadScreen(w fyne.Window, uploadFlow fyne.CanvasObject) *fyne.Container {
	return container.NewVBox(
		widget.NewButton("Upload Credit Card Transactions", func() {
			w.SetContent(uploadFlow)
		}),
		widget.NewButton("Upload Bank Transactions", func() {
			w.SetContent(uploadFlow)
		}),
	)
}
