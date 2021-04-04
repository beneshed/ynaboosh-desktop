package tables

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/models"
)

const (
	roundPad = 20
)

type internalTransactionTable struct {
	headers []string
	data    map[int]*models.Transaction
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
	if len(t.data) == 0 {
		return 1, len(t.headers)
	}
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
			if location.Row > 0 {
				row, _ := t.getHeaders()
				header := row[location.Col]
				if header == "Approved" && !t.data[location.Row].Approved {
					typedObj.Show()
				} else {
					typedObj.Hide()
				}
			} else {
				typedObj.Hide()
			}
		case *wrappedCheckedIcon:
			if location.Row > 0 {
				row, _ := t.getHeaders()
				header := row[location.Col]
				if header == "Approved" && t.data[location.Row].Approved {
					typedObj.Show()
				} else {
					typedObj.Hide()
				}
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
		data:    make(map[int]*models.Transaction),
		headers: []string{"Date of Transaction", "Payee", "Amount", "Category", "Approved", "Rule Detected"},
	}
	table := &TransactionTable{
		widget.NewTable(internal.Length, internal.createCell, internal.update),
		internal,
	}
	table.Table.OnSelected = table.OnSelected
	table.ExtendBaseWidget(table)
	return table
}

func (t *TransactionTable) OnSelected(id widget.TableCellID) {
	if len(t.internal.data) > 0 {
		if id.Row > 0 && id.Col == 4 {
			t.internal.data[id.Row].Approved = !t.internal.data[id.Row].Approved
			t.Refresh()
			return
		}
	}
}

func (t TransactionTable) GetHeaders() ([]string, error) {
	return t.internal.getHeaders()
}

func (t *TransactionTable) WrapTableWidth() {
	var maxDateWidth float32 = fyne.MeasureText(t.internal.headers[0], theme.TextSize(), fyne.TextStyle{}).Width
	var maxPayeeWidth float32 = fyne.MeasureText(t.internal.headers[1], theme.TextSize(), fyne.TextStyle{}).Width
	var maxAmountWidth float32 = fyne.MeasureText(t.internal.headers[2], theme.TextSize(), fyne.TextStyle{}).Width
	var maxCategoryWidth float32 = fyne.MeasureText(t.internal.headers[3], theme.TextSize(), fyne.TextStyle{}).Width
	var maxRuleDetectedWidth float32 = fyne.MeasureText(t.internal.headers[5], theme.TextSize(), fyne.TextStyle{}).Width
	var maxApprovedWidth float32 = fyne.MeasureText(t.internal.headers[4], theme.TextSize(), fyne.TextStyle{}).Width
	t.Table.SetColumnWidth(0, maxDateWidth+roundPad)
	t.Table.SetColumnWidth(1, maxPayeeWidth+roundPad)
	t.Table.SetColumnWidth(2, maxAmountWidth+roundPad)
	t.Table.SetColumnWidth(3, maxCategoryWidth+roundPad)
	t.Table.SetColumnWidth(4, maxApprovedWidth+roundPad)
	t.Table.SetColumnWidth(5, maxRuleDetectedWidth+roundPad)
	t.Table.Refresh()
}

func (t *TransactionTable) AddTransactions(transactions []models.Transaction) {
	var maxDateWidth float32 = fyne.MeasureText(t.internal.headers[0], theme.TextSize(), fyne.TextStyle{}).Width
	var maxPayeeWidth float32 = fyne.MeasureText(t.internal.headers[1], theme.TextSize(), fyne.TextStyle{}).Width
	var maxAmountWidth float32 = fyne.MeasureText(t.internal.headers[2], theme.TextSize(), fyne.TextStyle{}).Width
	var maxCategoryWidth float32 = fyne.MeasureText(t.internal.headers[3], theme.TextSize(), fyne.TextStyle{}).Width
	var maxRuleDetectedWidth float32 = fyne.MeasureText(t.internal.headers[5], theme.TextSize(), fyne.TextStyle{}).Width
	var maxApprovedWidth float32 = fyne.MeasureText(t.internal.headers[4], theme.TextSize(), fyne.TextStyle{}).Width
	lastRow := len(t.internal.data)
	for i, transaction := range transactions {
		tMap := transaction.GetTransactionAsMap()
		dateWidth := fyne.MeasureText(tMap["Date of Transaction"], theme.TextSize(), fyne.TextStyle{})
		payeeWidth := fyne.MeasureText(tMap["Payee"], theme.TextSize(), fyne.TextStyle{})
		amountWidth := fyne.MeasureText(tMap["Amount"], theme.TextSize(), fyne.TextStyle{})
		categoryWidth := fyne.MeasureText(tMap["Category"], theme.TextSize(), fyne.TextStyle{})
		ruleDetectedWidth := fyne.MeasureText(tMap["Rule Detected"], theme.TextSize(), fyne.TextStyle{})
		if dateWidth.Width > maxDateWidth {
			maxDateWidth = dateWidth.Width
		}
		if payeeWidth.Width > maxPayeeWidth {
			maxPayeeWidth = payeeWidth.Width
		}
		if amountWidth.Width > maxAmountWidth {
			maxAmountWidth = amountWidth.Width
		}
		if categoryWidth.Width > maxCategoryWidth {
			maxCategoryWidth = categoryWidth.Width
		}
		if ruleDetectedWidth.Width > maxRuleDetectedWidth {
			maxRuleDetectedWidth = ruleDetectedWidth.Width
		}
		transactionPtr := transaction
		t.internal.data[lastRow+i] = &transactionPtr
	}
	t.Table.SetColumnWidth(0, maxDateWidth+roundPad)
	t.Table.SetColumnWidth(1, maxPayeeWidth+roundPad)
	t.Table.SetColumnWidth(2, maxAmountWidth+roundPad)
	t.Table.SetColumnWidth(3, maxCategoryWidth+roundPad)
	t.Table.SetColumnWidth(4, maxApprovedWidth+roundPad)
	t.Table.SetColumnWidth(5, maxRuleDetectedWidth+roundPad)
	t.Table.Refresh()

}

func (t *TransactionTable) ApprovedSetAll(approved bool) {
	for _, transaction := range t.internal.data {
		transaction.Approved = approved
	}
	t.Table.Refresh()
}
