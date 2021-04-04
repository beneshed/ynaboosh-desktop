package institutions

import (
	"strconv"
	"strings"
	"time"

	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/models"
)

type Institution interface {
	ParseTransactions(fileName string) []models.Transaction
	Name() string
	DisplayName() string
}

func parseSlashedDate(date string, seperator *string, isInternational bool) time.Time {
	var dateParsed []string
	var day, month int
	if seperator == nil {
		dateParsed = strings.Split(date, "/")
	} else {
		dateParsed = strings.Split(date, *seperator)
	}
	if isInternational {
		day, _ = strconv.Atoi(strings.TrimSpace(dateParsed[0]))
		month, _ = strconv.Atoi(strings.TrimSpace(dateParsed[1]))
	} else {
		day, _ = strconv.Atoi(strings.TrimSpace(dateParsed[1]))
		month, _ = strconv.Atoi(strings.TrimSpace(dateParsed[0]))
	}

	year, _ := strconv.Atoi(strings.TrimSpace(dateParsed[2]))
	return time.Date(year+2000, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func LookupInstitutions() map[string]Institution {
	return map[string]Institution{
		discountBank: Discount{},
		leumiBank:    Leumi{},
		calCard:      Cal{},
		maxCard:      Max{},
	}
}
