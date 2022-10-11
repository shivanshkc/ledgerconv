package banks

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shivanshkc/ledgerconv/core/models"
)

// convHDFCSavings converts the HDFC savings account statements to JSON.
//
//nolint:funlen,cyclop // Converter functions can be long.
func convHDFCSavings(csvContent [][]string) ([]*models.TransactionDoc, error) {
	// Bank statement CSV files do not just contain the transaction list, but also some other metadata about the
	// bank account. This header allows us to detect the starting of the transaction table, so we can skip the needless.
	startingHeader := []string{
		"Date", "Narration", "Value Dat", "Debit Amount", "Credit Amount", "Chq/Ref Number", "Closing Balance",
	}

	// This var will hold the index of the first transaction table row.
	var startingIdx int
	// This var will hold the final list of converted transactions.
	var txDocs []*models.TransactionDoc //nolint:prealloc // Cannot pre-allocate this one.

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

	// Begin looping over the transaction table.
	for _, row := range csvContent[startingIdx:] {
		// Due to some reason, the statements contain empty rows in between too.
		if row[0] == "" {
			continue
		}

		// Parse timestamp.
		timestamp, err := time.Parse("02/01/06", row[0])
		if err != nil {
			// If we fail to parse the timestamp, we consider it as the end of the transaction table.
			return txDocs, nil //nolint:nilerr
		}

		// Other required fields.
		remarks, refNum := row[1], row[5]

		// Get the amount information.
		debitAmount, errDebit := strconv.ParseFloat(row[3], 64)
		creditAmount, errCredit := strconv.ParseFloat(row[4], 64)

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
		doc := &models.TransactionDoc{
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
