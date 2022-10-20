package converters

import (
	"strconv"
	"time"

	"github.com/shivanshkc/ledgerconv/core/models"
)

// ICICICredit is the converter function for ICICI credit card account statements.
//
//nolint:funlen // Converter functions are big. They parse bank statements!
func ICICICredit(csvContent [][]string) ([]*models.ConvertedTransactionDoc, error) {
	// Header of the transaction table in the CSV file.
	header := []string{
		"Date", "Sr.No.", "Transaction Details", "Reward Point Header",
		"Intl.Amount", "Amount(in Rs)", "BillingAmountSign",
	}

	// Trim each element. Bank statement schemas are not to be trusted!
	csvContent = trimCSV(csvContent)
	// startingIndex is the index of the first transaction row.
	startingIdx := getHeaderIndex(header, csvContent) + 2
	// Just a safety check.
	if startingIdx == 0 || startingIdx >= len(csvContent) {
		return nil, nil
	}

	// This var will hold the final list of converted transactions.
	var statement []*models.ConvertedTransactionDoc //nolint:prealloc // Cannot pre-alloc.
	// Begin looping over the transaction table.
	for _, row := range csvContent[startingIdx:] {
		// Due to some unearthly reason, the statements contain empty rows in between too.
		if len(row) == 0 || row[0] == "" {
			continue
		}

		// Parse timestamp.
		timestamp, err := time.Parse("02/01/2006", row[0])
		if err != nil {
			// If we fail to parse the timestamp, we consider it as the end of the transaction table.
			return statement, nil //nolint:nilerr
		}

		// Other required fields.
		serial, remarks, amountSign := row[1], row[2], row[6]

		// Get the amount information.
		amount, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			// If we fail to parse the amount, we consider it as the end of the transaction table.
			return statement, nil //nolint:nilerr
		}

		// If the amount is zero, we do not consider this row.
		if amount == 0 {
			continue
		}

		// Apparently, debit transactions have "CR" BillingAmountSign -_-
		if amountSign == "CR" {
			amount *= -1
		}

		// Instantiating the transaction doc.
		doc := &models.ConvertedTransactionDoc{
			AccountName:     "", // This is not the responsibility of the converterFunc.
			Amount:          amount,
			Timestamp:       timestamp,
			BankSerial:      serial,
			BankPaymentMode: "",
			BankRemarks:     remarks,
		}

		// Collecting the result.
		statement = append(statement, doc)
	}

	return statement, nil
}
