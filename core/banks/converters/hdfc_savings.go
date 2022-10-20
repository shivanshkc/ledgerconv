package converters

import (
	"strconv"
	"time"

	"github.com/shivanshkc/ledgerconv/core/models"
)

// HDFCSavings is the converter function for HDFC savings account statements.
func HDFCSavings(csvContent [][]string) ([]*models.ConvertedTransactionDoc, error) {
	// Header of the transaction table in the CSV file.
	header := []string{
		"Date", "Narration", "Value Dat", "Debit Amount", "Credit Amount", "Chq/Ref Number", "Closing Balance",
	}

	return base(csvContent, header, 1, hdfcSavingsRow)
}

// hdfcSavingsRow is the rowConverterFunc for HDFC savings account statements.
func hdfcSavingsRow(row []string) (*models.ConvertedTransactionDoc, bool, error) {
	// Due to some reason, the statements contain empty rows in between too.
	if row[0] == "" {
		// Skipping to the next row.
		return nil, false, nil
	}

	// Parse timestamp.
	timestamp, err := time.Parse("02/01/06", row[0])
	if err != nil {
		// If we fail to parse the timestamp, we consider it as the end of the transaction table.
		return nil, true, nil //nolint:nilerr
	}

	// Other required fields.
	remarks, serial := row[1], row[5]

	// Get the amount information.
	debitAmount, errDebit := strconv.ParseFloat(row[3], 64)
	creditAmount, errCredit := strconv.ParseFloat(row[4], 64)

	// If both amounts failed to parse, we cannot proceed further.
	if errDebit != nil && errCredit != nil {
		// If we fail to parse the amounts, we consider it as the end of the transaction table.
		return nil, true, nil //nolint:nilerr
	}

	// Debit amount is taken negative.
	amount := -debitAmount
	// If debit amount is invalid or zero, we take the credit amount as the final value.
	if errDebit != nil || debitAmount == 0 {
		amount = creditAmount
	}

	// If the amount is zero, we do not consider this row.
	if amount == 0 {
		// Skipping to the next row.
		return nil, false, nil
	}

	// Instantiating the transaction doc.
	return &models.ConvertedTransactionDoc{
		AccountName:     "", // This is not the responsibility of the converterFunc.
		Amount:          amount,
		Timestamp:       timestamp,
		BankSerial:      serial,
		BankPaymentMode: "",
		BankRemarks:     remarks,
	}, false, nil
}
