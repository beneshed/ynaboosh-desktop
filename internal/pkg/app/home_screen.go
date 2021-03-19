package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func NewHomeScreen(fileUploadFlow Component, settingsScreen Component) *fyne.Container {
	tabs := container.NewVBox(container.NewAppTabs(
		container.NewTabItem("File Uploads", fileUploadFlow.Object),
		//container.NewTabItem("History", widget.NewLabel("Historical Uploaded Transactions")),
		//container.NewTabItem("Suggested Categories", table),
		container.NewTabItem("Settings", settingsScreen.Object),
	))
	return tabs
}
