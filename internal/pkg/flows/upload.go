package flows

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/thebenwaters/ynab-desktop-importer/internal/pkg/ynabimporter"
)

type InternalState ynabimporter.GlobalState

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

func fetchAccountsForSelector(accounts []ynabimporter.Account) []string {
	var results []string
	for _, account := range accounts {
		if !account.Deleted && !account.Closed {
			log.Println(account.Type)
			results = append(results, account.Name)
		}
	}
	return results
}

func (s InternalState) NewUploadFlow(w fyne.Window) *fyne.Container {
	state := NewUploadState()
	fileSelected := widget.NewLabel("")
	fileLabel := widget.NewButton("File to Upload", func() {
		dialog.ShowFileOpen(func(info fyne.URIReadCloser, err error) {
			fileSelected.SetText(info.URI().Path())
			fileToUpload := info.URI().Path()
			fileToUploadPtr := &fileToUpload
			state.fileToUpload = fileToUploadPtr
		}, w)

	})
	accountSelectLabel := widget.NewLabel("YNAB Account")
	var accountsToQuery []ynabimporter.Account
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
	return fyne.NewContainer(container.New(layout.NewFormLayout(), fileLabel, fileSelected, accountSelectLabel, accountSelector, cancelButton, nextButton))
}
