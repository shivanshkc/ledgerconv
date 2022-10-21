package converters

import (
	"strconv"
	"time"

	"github.com/shivanshkc/ledgerconv/core/models"
)

// ICICISavings is the converter function for ICICI savings account statements.
//
//nolint:funlen,cyclop // Converter functions are big. They parse bank statements!
func ICICISavings(csvContent [][]string) ([]*models.ConvertedTransactionDoc, error) {
	// Header of the transaction table in the CSV file.
	header := []string{"DATE", "MODE", "PARTICULARS", "DEPOSITS", "WITHDRAWALS", "BALANCE"}

	// Trim each element. Bank statement schemas are not to be trusted!
	csvContent = trimCSV(csvContent)
	// startingIndex is the index of the first transaction row.
	startingIdx := getHeaderIndex(header, csvContent) + 1
	// Just a safety check.
	if startingIdx == 0 || startingIdx >= len(csvContent) {
		return nil, nil
	}

	// This var will hold the final list of converted transactions.
	var statement []*models.ConvertedTransactionDoc //nolint:prealloc // Cannot pre-alloc.
	// Begin looping over the transaction table.
	for _, row := range csvContent[startingIdx+1:] {
		// Due to some unearthly reason, the statements contain empty rows in between too.
		if row[0] == "" {
			continue
		}

		// Parse timestamp.
		timestamp, err := time.Parse("02-01-2006", row[0])
		if err != nil {
			// If we fail to parse the timestamp, we consider it as the end of the transaction table.
			return statement, nil //nolint:nilerr
		}

		// Other required fields.
		paymentMode, remarks := row[1], row[2]

		// Get the amount information.
		creditAmount, errCredit := strconv.ParseFloat(row[3], 64)
		debitAmount, errDebit := strconv.ParseFloat(row[4], 64)

		// If both amounts failed to parse, we cannot proceed further.
		if errDebit != nil && errCredit != nil {
			// If we fail to parse the amounts, we consider it as the end of the transaction table.
			return statement, nil //nolint:nilerr
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

		// Instantiate the transaction doc.
		doc := &models.ConvertedTransactionDoc{
			AccountName:     "", // This is not the responsibility of the converterFunc.
			Amount:          amount,
			Timestamp:       timestamp,
			BankSerial:      "",
			BankPaymentMode: paymentMode,
			BankRemarks:     remarks,
		}

		// Collect the result.
		statement = append(statement, doc)
	}

	return statement, nil
}
