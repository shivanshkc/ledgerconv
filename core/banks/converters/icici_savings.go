package converters

import (
	"strconv"
	"time"

	"github.com/shivanshkc/ledgerconv/core/models"
)

// ICICISavings is the converter function for ICICI savings account statements.
func ICICISavings(csvContent [][]string) ([]*models.ConvertedTransactionDoc, error) {
	header := []string{"DATE", "MODE", "PARTICULARS", "DEPOSITS", "WITHDRAWALS", "BALANCE"}
	return base(csvContent, header, 1, iciciSavingsRow)
}

// iciciSavingsRow is the rowConverterFunc for ICICI savings account statements.
func iciciSavingsRow(row []string) (*models.ConvertedTransactionDoc, bool, error) {
	// Due to some unearthly reason, the statements contain empty rows in between too.
	if row[0] == "" {
		// Skip to the next row.
		return nil, false, nil
	}

	// Parse timestamp.
	timestamp, err := time.Parse("02-01-2006", row[0])
	if err != nil {
		// If we fail to parse the timestamp, we consider it as the end of the transaction table.
		return nil, true, nil //nolint:nilerr
	}

	// Other required fields.
	paymentMode, remarks := row[1], row[2]

	// Get the amount information.
	creditAmount, errCredit := strconv.ParseFloat(row[3], 64)
	debitAmount, errDebit := strconv.ParseFloat(row[4], 64)

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
		// Skip to the next row.
		return nil, false, nil
	}

	// Instantiating the transaction doc.
	return &models.ConvertedTransactionDoc{
		AccountName:     "", // This is not the responsibility of the converterFunc.
		Amount:          amount,
		Timestamp:       timestamp,
		BankSerial:      "",
		BankPaymentMode: paymentMode,
		BankRemarks:     remarks,
	}, false, nil
}
