package app

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type UploadState struct {
	fileToUpload *string
	ynabAccount  *string
	uploadType   string
}

func NewUploadState() *UploadState {
	return &UploadState{}
}

func (s UploadState) Validate() bool {
	if s.fileToUpload == nil || s.ynabAccount == nil {
		return false
	}
	return true
}

func fetchAccountsForSelector(accounts []Account) []string {
	var results []string
	for _, account := range accounts {
		if !account.Deleted && !account.Closed {
			results = append(results, account.Name)
		}
	}
	return results
}

func (s GlobalState) NewUploadFlow() *Component {
	state := NewUploadState()
	fileSelected := widget.NewLabel("")
	fileLabel := widget.NewButton("File to Upload", func() {
		dialog.ShowFileOpen(func(info fyne.URIReadCloser, err error) {
			if err != nil && info != nil {
				fileSelected.SetText(info.URI().Path())
				fileToUpload := info.URI().Path()
				fileToUploadPtr := &fileToUpload
				state.fileToUpload = fileToUploadPtr
			}
		}, s.Window)
	})
	accountSelectLabel := widget.NewLabel("YNAB Account")
	var accountsToQuery []Account
	log.Println(s.DB)
	results := s.DB.Find(&accountsToQuery)
	if results.Error != nil {
		log.Println(results.Error)
	}
	accountSelector := widget.NewSelect(fetchAccountsForSelector(accountsToQuery), func(option string) {
		optionPtr := &option
		state.ynabAccount = optionPtr
	})
	cancelButton := widget.NewButton("Cancel", func() {
		log.Println("Clicked")
	})
	nextButton := widget.NewButton("Next", func() {
		log.Println(state.Validate())
	})
	return NewComponent(container.NewVBox(
		widget.NewSeparator(),
		container.New(layout.NewFormLayout(), fileLabel, fileSelected, accountSelectLabel, accountSelector, cancelButton, nextButton)))
}
