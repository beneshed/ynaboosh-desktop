package institutions

import (
	"log"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/language"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/models"
)

const (
	discountBankDisplayName = "Discount - Bank"
	discountBank            = "DISCOUNT"
)

type Discount struct{}

func (i Discount) Name() string {
	return discountBank
}

func (i Discount) DisplayName() string {
	return discountBankDisplayName
}

func (i Discount) ParseTransactions(fileName string) []models.Transaction {
	var transactions []models.Transaction
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Panicln(err)
	}
	sheets := f.GetSheetList()
	sheet := sheets[0]
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
			amountDecimalWrapped := float32(amountDecimal)
			if err != nil {
				log.Panicln(err)
			}
			payee := row[2]
			reversedPayee := ""
			payee = language.ReverseHebrew(payee)
			postiveAmount := amountDecimalWrapped
			if amountDecimalWrapped < 0.0 {
				postiveAmount = amountDecimalWrapped * -1.0
			}
			transactions = append(transactions, models.Transaction{
				DateOfTransaction: parseSlashedDate(row[0], &weirdDate, false),
				Payee:             payee,
				ReversePayee:      reversedPayee,
				Amount:            postiveAmount,
				Out:               amountDecimalWrapped < 0.0,
			})
		}
	}
	return transactions
}
