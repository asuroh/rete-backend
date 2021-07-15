package helper

import (
	"strconv"
	"time"
)

// GenerateInvoice ...
func GenerateInvoice(count int) string {
	now := time.Now().Format("06-01-02")

	invoice := now + "00001"

	if count < 9 {
		invoice = SplitDate(now) + "0000" + strconv.Itoa(count+1)
	} else if count < 99 {
		invoice = SplitDate(now) + "000" + strconv.Itoa(count+1)
	} else if count < 999 {
		invoice = SplitDate(now) + "00" + strconv.Itoa(count+1)
	} else if count < 9999 {
		invoice = SplitDate(now) + "0" + strconv.Itoa(count+1)
	} else {
		invoice = SplitDate(now) + strconv.Itoa(count+1)
	}

	return "INV" + invoice
}
