package core

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

// convICICICredit converts the ICICI credit card statements to JSON.
//
//nolint:funlen // Converter functions can be long.
func convICICICredit(csvContent [][]string) ([]*transactionDoc, error) {
	// Bank statement CSV files do not just contain the transaction list, but also some other metadata about the
	// bank account. This header allows us to detect the starting of the transaction table, so we can skip the needless.
	startingHeader := []string{
		"Date", "Sr.No.", "Transaction Details", "Reward Point Header", "Intl.Amount",
		"Amount(in Rs)", "BillingAmountSign",
	}

	// This var will hold the index of the first transaction table row.
	var startingIdx int
	// This var will hold the final list of converted transactions.
	var txDocs []*transactionDoc //nolint:prealloc // Cannot pre-allocate this one.

	// Loop over CSV rows to find the starting of the transaction table.
	//nolint:varnamelen // "i" is a fine name here.
	for i, row := range csvContent {
		// If the row is not the same as the starting header, we continue.
		if !reflect.DeepEqual(row, startingHeader) {
			continue
		}
		// Starting of the transaction table is located.
		startingIdx = i + 2
		break
	}

	// Just a safety check.
	if startingIdx >= len(csvContent) {
		return nil, nil
	}

	// Begin looping over the transaction table.
	for _, row := range csvContent[startingIdx:] {
		// Trim each element. Bank statement schemas are not to be trusted!
		for i := range row {
			row[i] = strings.TrimSpace(row[i])
		}

		// Parse timestamp.
		timestamp, err := time.Parse("02/01/2006", row[0])
		if err != nil {
			// If we fail to parse the timestamp, we consider it as the end of the transaction table.
			return txDocs, nil //nolint:nilerr
		}

		// Other required fields.
		refNum, remarks, amountSign := row[1], row[2], row[6]

		// Get the amount information.
		amount, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			// If we fail to parse the amount, we consider it as the end of the transaction table.
			return txDocs, nil //nolint:nilerr
		}

		// If the amount is zero, we do not consider this row.
		if amount == 0 {
			continue
		}

		// If the amount sign is CR, it means it is a credit transaction.
		if amountSign != "CR" {
			amount *= -1
		}

		// Instantiating the transaction doc.
		doc := &transactionDoc{
			AccountName: "", // This is not the responsibility of the converterFunc.
			Amount:      amount,
			Timestamp:   timestamp,
			RefNum:      refNum,
			PaymentMode: "",
			Remarks:     remarks,
		}

		// Collecting the result.
		txDocs = append(txDocs, doc)
	}

	return txDocs, nil
}
