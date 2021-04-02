package utils

import (
	"errors"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/thebenwaters/ynab-utils-desktop/institutions"
)

type internalTransactionTable struct {
	headers []string
	data    map[int]*institutions.Transaction
}

type wrappedUncheckedIcon struct {
	*widget.Icon
}

type wrappedCheckedIcon struct {
	*widget.Icon
}

func newWrappedUncheckedIcon() *wrappedUncheckedIcon {
	icon := &wrappedUncheckedIcon{
		widget.NewIcon(theme.CheckButtonIcon()),
	}
	icon.ExtendBaseWidget(icon)
	return icon
}

func newWrappedCheckedIcon() *wrappedCheckedIcon {
	icon := &wrappedCheckedIcon{
		widget.NewIcon(theme.CheckButtonCheckedIcon()),
	}
	icon.ExtendBaseWidget(icon)
	return icon
}

func (t *internalTransactionTable) Length() (int, int) {
	return len(t.data), len(t.headers)
}

func (t internalTransactionTable) getHeaders() ([]string, error) {
	if len(t.headers) < 1 || t.headers == nil {
		return []string{}, errors.New("headers not initialized")
	}
	return t.headers, nil
}

func (t *internalTransactionTable) createCell() fyne.CanvasObject {
	uncheckedIcon := newWrappedUncheckedIcon()
	uncheckedIcon.Hide()
	checkedIcon := newWrappedCheckedIcon()
	checkedIcon.Hide()
	return container.NewMax(widget.NewLabel(""), uncheckedIcon, checkedIcon) //, button)
}

func (t *internalTransactionTable) update(location widget.TableCellID, cell fyne.CanvasObject) {
	c := cell.(*fyne.Container)
	for _, obj := range c.Objects {
		switch typedObj := obj.(type) {
		case *widget.Label:
			if location.Row == 0 {
				row, _ := t.getHeaders()
				headers := row[location.Col]
				typedObj.SetText(headers)
			} else {
				row := t.data[location.Row]
				rowMap := row.GetTransactionAsMap()
				rowOfHeaders, _ := t.getHeaders()
				value := rowMap[rowOfHeaders[location.Col]]
				typedObj.SetText(value)
			}
		case *wrappedUncheckedIcon:
			row, _ := t.getHeaders()
			header := row[location.Col]
			if header == "Approved" && !t.data[location.Row].Approved && location.Row > 0 {
				typedObj.Show()
			} else {
				typedObj.Hide()
			}
		case *wrappedCheckedIcon:
			row, _ := t.getHeaders()
			header := row[location.Col]
			if header == "Approved" && t.data[location.Row].Approved && location.Row > 0 {
				typedObj.Show()
			} else {
				typedObj.Hide()
			}
		}
	}
}

type TransactionTable struct {
	*widget.Table
	internal internalTransactionTable
}

func NewTransactionTable() *TransactionTable {

	internal := internalTransactionTable{
		data:    make(map[int]*institutions.Transaction),
		headers: []string{"Date of Transaction", "Payee", "Amount", "Approved"},
	}
	table := &TransactionTable{
		widget.NewTable(internal.Length, internal.createCell, internal.update),
		internal,
	}
	for i := range internal.headers {
		table.Table.SetColumnWidth(i, 350)
	}
	table.Table.OnSelected = table.OnSelected
	table.ExtendBaseWidget(table)
	return table
}

func (t *TransactionTable) OnSelected(id widget.TableCellID) {
	headers, err := t.internal.getHeaders()
	if err != nil {
		log.Panicln("headers shit")
	}
	if id.Row > 0 && id.Col == len(headers)-1 {
		t.internal.data[id.Row].Approved = !t.internal.data[id.Row].Approved
		t.Refresh()
		return
	}
}

func (t *TransactionTable) AddTransactions(transactions []institutions.Transaction) {
	lastRow := len(t.internal.data)
	for i, transaction := range transactions {
		transactionPtr := transaction
		t.internal.data[lastRow+i] = &transactionPtr
	}
}

func (t *TransactionTable) ApprovedSetAll(approved bool) {
	for _, transaction := range t.internal.data {
		transaction.Approved = approved
	}
}
