package banks

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shivanshkc/ledgerconv/core/models"
)

// convICICISavings converts the ICICI savings account statements to JSON.
//
//nolint:funlen,cyclop // Converter functions can be long.
func convICICISavings(csvContent [][]string) ([]*models.ConvertedTransactionDoc, error) {
	// Bank statement CSV files do not just contain the transaction list, but also some other metadata about the
	// bank account. This header allows us to detect the starting of the transaction table, so we can skip the needless.
	startingHeader := []string{"DATE", "MODE", "PARTICULARS", "DEPOSITS", "WITHDRAWALS", "BALANCE"}

	// This var will hold the index of the first transaction table row.
	var startingIdx int
	// This var will hold the final list of converted transactions.
	var txDocs []*models.ConvertedTransactionDoc //nolint:prealloc // Cannot pre-allocate this one.

	// Trim each element. Bank statement schemas are not to be trusted!
	for i := range csvContent {
		for j := range csvContent[i] {
			csvContent[i][j] = strings.TrimSpace(csvContent[i][j])
		}
	}

	// Loop over CSV rows to find the starting of the transaction table.
	//nolint:varnamelen // "i" is a fine name here.
	for i, row := range csvContent {
		// If the row is not the same as the starting header, we continue.
		if !reflect.DeepEqual(row, startingHeader) {
			continue
		}
		// Starting of the transaction table is located.
		startingIdx = i + 1
		break
	}

	// Just a safety check.
	if startingIdx == 0 || startingIdx >= len(csvContent) {
		return nil, nil
	}

	for _, row := range csvContent[startingIdx+1:] {
		// Due to some reason, the statements contain empty rows in between too.
		if row[0] == "" {
			continue
		}

		// Parse timestamp.
		timestamp, err := time.Parse("02-01-2006", row[0])
		if err != nil {
			// If we fail to parse the timestamp, we consider it as the end of the transaction table.
			return txDocs, nil //nolint:nilerr
		}

		// Other required fields.
		paymentMode, remarks := row[1], row[2]

		// Get the amount information.
		creditAmount, errCredit := strconv.ParseFloat(row[3], 64)
		debitAmount, errDebit := strconv.ParseFloat(row[4], 64)

		// If both amounts failed to parse, we cannot proceed further.
		if errDebit != nil && errCredit != nil {
			// If we fail to parse the amounts, we consider it as the end of the transaction table.
			return txDocs, nil //nolint:nilerr
		}

		// Debit amount is taken negative.
		amount := -debitAmount
		// If debit amount is invalid or zero, we take the credit amount as the final value.
		if errDebit != nil || debitAmount == 0 {
			amount = creditAmount
		}

		// If the amount is zero, we do not consider this row.
		if amount == 0 {
			continue
		}

		// Instantiating the transaction doc.
		doc := &models.ConvertedTransactionDoc{
			AccountName: "", // This is not the responsibility of the converterFunc.
			Amount:      amount,
			Timestamp:   timestamp,
			RefNum:      "",
			PaymentMode: paymentMode,
			Remarks:     remarks,
		}

		// Collecting the result.
		txDocs = append(txDocs, doc)
	}

	return txDocs, nil
}
