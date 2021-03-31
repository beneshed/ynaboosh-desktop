package main

import (
	"errors"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/thebenwaters/ynab-utils-desktop/institutions"
)

const (
	tagName = "ynab"
)

type internalTransactionTable struct {
	data map[int]interface{}
}

func (t *internalTransactionTable) length() (int, int) {
	headers, ok := t.data[0].([]string)
	if !ok {
		return 0, 0
	}
	return len(t.data), len(headers)
}

func (t internalTransactionTable) headers() ([]string, error) {
	headers, ok := t.data[0].([]string)
	if !ok {
		return []string{}, errors.New("Headers Not Initialized")
	}
	return headers, nil
}

func (t internalTransactionTable) getDataRow(row int) (*institutions.Transaction, error) {
	data, ok := t.data[row].(institutions.Transaction)
	if !ok {
		return &institutions.Transaction{}, errors.New("Not found")
	}
	return &data, nil
}

func (t *internalTransactionTable) createCell() fyne.CanvasObject {
	button := widget.NewCheck("", func(c bool) {
		log.Println("Clicked")
	})
	button.Hide()
	return container.NewMax(widget.NewLabel(""), button)
}

func (t *internalTransactionTable) update(location widget.TableCellID, cell fyne.CanvasObject) {
	c := cell.(*fyne.Container)
	for _, obj := range c.Objects {
		switch typedObj := obj.(type) {
		case *widget.Label:
			if location.Row == 0 {
				row, _ := t.headers()
				headers := row[location.Col]
				typedObj.SetText(headers)
			} else {
				row, _ := t.getDataRow(location.Row)
				rowMap := row.GetTransactionAsMap()
				rowOfHeaders, _ := t.headers()
				header := rowOfHeaders[location.Col]
				value := rowMap[header]
				typedObj.SetText(value)
			}
		case *widget.Check:
			row, _ := t.getDataRow(location.Row)
			rowMap := row.GetTransactionAsMap()
			rowOfHeaders, _ := t.headers()
			header := rowOfHeaders[location.Col]
			value := rowMap[header]
			if location.Row > 0 && location.Col == len(rowOfHeaders)-1 {
				log.Println("HERE", location.Row, location.Col, value)
				typedObj.Show()
			}
		}
	}
}

type TransactionTable struct {
	*widget.Table
	internal internalTransactionTable
}

func NewTransactionTable() *TransactionTable {
	headers := []string{"Date of Transaction", "Payee", "Amount", "Approved"}
	internal := internalTransactionTable{
		data: map[int]interface{}{
			0: headers,
		},
	}
	table := &TransactionTable{
		widget.NewTable(internal.length, internal.createCell, internal.update),
		internal,
	}
	for i, _ := range headers {
		table.Table.SetColumnWidth(i, 350)
	}
	table.Table.OnSelected = table.OnSelected
	table.ExtendBaseWidget(table)
	return table
}

func (t *TransactionTable) OnSelected(id widget.TableCellID) {
	headers, err := t.internal.headers()
	if err != nil {
		log.Panicln("shit")
	}
	if id.Row > 0 && id.Col == len(headers)-1 {
		row, _ := t.internal.getDataRow(id.Row)
		log.Println(row)
		row.Approved = !row.Approved
		log.Println(row)
		return
		//log.Println("SELECTED APPROVED", t, id)
	}
}

func (t *TransactionTable) AddTransactions(transactions []institutions.Transaction) {
	lastRow := len(t.internal.data)
	for i, transaction := range transactions {
		t.internal.data[lastRow+i] = transaction
	}
}

func (t *TransactionTable) UpdateTransactions(tranactions map[int]institutions.Transaction) {
	for k, v := range transactions {
		t.internal.data[k] = v
	}
}
