package institutions

import (
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
)

const (
	leumiTableXPath = `//*[@id="ctlActivityTable"]//tr//td`
	// `//*[@id="ctlActivityTable"]/tbody/tr[1]/th[1]/a/span`
)

func ParseLeumiTransactions(fileName string) []Transaction {
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panicln(err)
	}
	doc, err := htmlquery.Parse(strings.NewReader(string(fileBytes)))
	if err != nil {
		log.Panicln(err)
	}
	log.Println(doc)
	results, err := htmlquery.QueryAll(doc, leumiTableXPath)
	if err != nil {
		log.Println(err)
	}
	var transactions []Transaction
	for i := 0; i < len(results); i += 6 {
		date := parseSlashedDate(htmlquery.InnerText(results[i]), nil)
		payee := strings.TrimSpace(htmlquery.InnerText(results[i+1]))
		if detectHebrew(payee) {
			payee = reverse(payee)
		}
		isOut := false
		out, _ := strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(htmlquery.InnerText(results[i+3])), ",", ""), 32)
		in, _ := strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(htmlquery.InnerText(results[i+4])), ",", ""), 32)
		if out > 0 {
			isOut = true
		}
		transactions = append(transactions, Transaction{
			DateOfTransaction: date,
			Payee:             payee,
			CurrencyCode:      nis,
			Amount:            float32(math.Max(out, in)),
			Out:               isOut,
		})
	}
	return transactions
}
