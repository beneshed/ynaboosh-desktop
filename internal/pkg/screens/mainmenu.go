package screens

import (
	"log"

	"fyne.io/fyne/v2"
)

func NewMainMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("foo", func() {
				log.Println("foo")
			}),
			fyne.NewMenuItem("bar", func() {
				log.Println("bar")
			})),
	)
}
