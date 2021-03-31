package institutions

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func parseMaxDate(date string) time.Time {
	dateParsed := strings.Split(date, "-")
	day, _ := strconv.Atoi(strings.TrimSpace(dateParsed[0]))
	month, _ := strconv.Atoi(strings.TrimSpace(dateParsed[1]))
	year, _ := strconv.Atoi(strings.TrimSpace(dateParsed[2]))
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func parseSheet(f *excelize.File, sheetName string) []Transaction {
	var transactions []Transaction
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Println(err)
	}
	headers := rows[3]
	for _, row := range rows[4:] {
		if len(row) == len(headers) {
			adjustedAmount := strings.ReplaceAll(row[5], ",", "")
			amountDecimal, err := strconv.ParseFloat(adjustedAmount, 32)
			amountDecimalWrapped := float32(amountDecimal)
			if err != nil {
				log.Panicln(err)
			}
			payee := row[1]
			if detectHebrew(payee) {
				payee = reverse(payee)
			}
			postiveAmount := amountDecimalWrapped
			if amountDecimalWrapped < 0.0 {
				postiveAmount = amountDecimalWrapped * -1.0
			}
			transactions = append(transactions, Transaction{
				DateOfTransaction: parseMaxDate(row[0]),
				Payee:             payee,
				Amount:            postiveAmount,
				Out:               amountDecimalWrapped < 0.0,
			})
		}
	}
	return transactions
}

func ParseMaxTransations(fileName string) []Transaction {
	var transactions []Transaction
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Panicln(err)
	}
	sheets := f.GetSheetList()
	shekelSheet := sheets[0]
	foreignSheet := sheets[2]
	shekelTransactions := parseSheet(f, shekelSheet)
	foreignTransactions := parseSheet(f, foreignSheet)
	transactions = append(transactions, shekelTransactions...)
	transactions = append(transactions, foreignTransactions...)
	return transactions
}
