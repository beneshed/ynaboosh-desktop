package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func NewHomeScreen(fileUploadFlow fyne.CanvasObject, settingsScreen fyne.CanvasObject, aboutScreen fyne.CanvasObject) *container.AppTabs {
	tabs := container.NewAppTabs(
		container.NewTabItem("File Uploads", fileUploadFlow),
		//container.NewTabItem("History", widget.NewLabel("Historical Uploaded Transactions")),
		//container.NewTabItem("Suggested Categories", table),
		container.NewTabItem("Settings", settingsScreen),
		container.NewTabItem("About", aboutScreen),
	)
	return tabs
}
