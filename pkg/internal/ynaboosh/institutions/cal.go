package institutions

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"

	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/models"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const (
	calCard            = "CAL"
	calCardDisplayName = "Cal - Credit Card"
)

type Cal struct {
	name string
}

func (i Cal) Name() string {
	return calCard
}

func (i Cal) DisplayName() string {
	return calCardDisplayName
}

func detectHebrew(s string) bool {
	words := regexp.MustCompile("[\u0590-\u05fe]+")
	return len(words.FindAllString(s, -1)) > 0
}

func detectShekel(s string) bool {
	shekel := regexp.MustCompile("₪")
	return len(shekel.FindAllString(s, -1)) > 0
}

func grabAmount(s string) float32 {
	amount := regexp.MustCompile(`\d+\.\d+`)
	matches := amount.FindAllString(s, -1)
	value, err := strconv.ParseFloat(matches[0], 32)
	if err != nil {
		log.Panicln(err)
	}
	return float32(value)
}

func reverse(s string) string {
	rns := []rune(s) // convert to rune
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {

		// swap the letters of the string,
		// like first with last and so on.
		rns[i], rns[j] = rns[j], rns[i]
	}

	// return the reversed string.
	return string(rns)
}

func (i Cal) ParseTransactions(fileName string) []models.Transaction {
	filesBytes, _ := ioutil.ReadFile(fileName)
	dec := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	utf16r := bytes.NewReader(filesBytes)
	utf8r := transform.NewReader(utf16r, dec)
	csvr := csv.NewReader(utf8r)
	csvr.Comma = '\t'
	csvr.FieldsPerRecord = -1
	csvr.LazyQuotes = true
	r, err := csvr.ReadAll()
	if err != nil {
		log.Panicln(err)
	}
	headers := r[2]
	data := r[3 : len(r)-1]
	var transactions []models.Transaction
	headersMap := make(map[int]string)
	for index, h := range headers {
		headersMap[index] = reverse(h)
	}
	for _, entry := range data {
		transaction := models.Transaction{}
		for index, col := range entry {
			switch index {
			case 0:
				transaction.DateOfTransaction = parseSlashedDate(col, nil, true)
			case 1:
				payee := col
				if detectHebrew(col) {
					transaction.ReversePayee = col
					payee = reverse(col)
				}
				transaction.Payee = payee
			case 3:
				value := grabAmount(col)
				if detectShekel(col) {
					transaction.CurrencyCode = "nis"
				} else {
					transaction.CurrencyCode = "usd"
				}
				transaction.Amount = value
			}
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}
