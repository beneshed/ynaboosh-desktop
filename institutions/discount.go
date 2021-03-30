package institutions

import (
	"log"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func ParseDiscountTransactions(fileName string, isHebrew bool) []Transaction {
	var transactions []Transaction
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Panicln(err)
	}
	sheets := f.GetSheetList()
	log.Println("Sheets", sheets)
	sheet := sheets[0]
	log.Println(sheet)
	rows, err := f.GetRows(sheet)
	if err != nil {
		log.Println(err)
	}
	headers := rows[12]
	weirdDate := "-"
	for _, row := range rows[13:] {
		if len(row) == len(headers) {
			adjustedAmount := strings.ReplaceAll(row[3], ",", "")
			amountDecimal, err := strconv.ParseFloat(adjustedAmount, 32)
			if err != nil {
				log.Panicln(err)
			}
			log.Println(row[0])
			payee := row[2]
			if detectHebrew(payee) {
				payee = reverse(payee)
			}
			transactions = append(transactions, Transaction{
				DateOfTransaction: parseSlashedDate(row[0], &weirdDate),
				Payee:             payee,
				Amount:            float32(amountDecimal) * -1.0,
				Out:               float32(amountDecimal) > 0.0,
			})
		}
	}
	return transactions
}
