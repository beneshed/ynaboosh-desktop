package forms

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/extensions"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/institutions"
)

func NewSyncTransactionsUploadForm(table *extensions.TransactionTable, window fyne.Window) *extensions.ClearableForm {
	var supportedFileTypes []string
	fileType := binding.NewString()
	fileSelected := binding.NewString()
	institutionsMap := institutions.LookupInstitutions()
	for _, i := range institutionsMap {
		supportedFileTypes = append(supportedFileTypes, i.DisplayName())
	}
	filePicker := dialog.NewFileOpen(func(info fyne.URIReadCloser, err error) {
		log.Println(info, err)
		if err == nil && info != nil {
			log.Println("should work")
			err = fileSelected.Set(info.URI().Path())
			log.Println(err)
		}
	}, window)
	// sync transactions
	// 1) form
	// 2) adjust / fix / confirm
	form := widget.NewForm(
		widget.NewFormItem(
			"", widget.NewButtonWithIcon("Transaction File", theme.FileIcon(), func() {
				filePicker.Show()
			}),
		),
		widget.NewFormItem(
			"File Type", widget.NewSelect(supportedFileTypes, func(value string) {
				for _, i := range institutionsMap {
					if i.DisplayName() == value {
						fileType.Set(i.Name())
					}
				}
			}),
		),
		widget.NewFormItem(
			"File Path", widget.NewLabelWithData(fileSelected),
		),
		widget.NewFormItem(
			"Account", widget.NewSelect([]string{}, func(value string) {
				log.Println(value)
			}),
		),
	)
	clearableForm := extensions.NewClearableForm(form)
	clearableForm.OnSubmit = func() {
		// run file detection
		fileName, err := fileSelected.Get()
		if err != nil {
			log.Panicln(err)
		}
		log.Println("About to parse")
		fileTypeLookup, err := fileType.Get()
		if err != nil {
			log.Panicln(err)
		}
		institution := institutionsMap[fileTypeLookup]
		transactions := institution.ParseTransactions(fileName)
		table.AddTransactions(transactions)
	}
	clearableForm.SubmitText = "Load File"
	return clearableForm
}
