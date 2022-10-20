package converters

import (
	"strconv"
	"time"

	"github.com/shivanshkc/ledgerconv/core/models"
)

// ICICICredit is the converter function for ICICI credit card account statements.
func ICICICredit(csvContent [][]string) ([]*models.ConvertedTransactionDoc, error) {
	// Header of the transaction table in the CSV file.
	header := []string{
		"Date", "Sr.No.", "Transaction Details", "Reward Point Header",
		"Intl.Amount", "Amount(in Rs)", "BillingAmountSign",
	}

	return base(csvContent, header, 2, iciciCreditRow)
}

// iciciCreditRow is the rowConverterFunc for ICICI credit card account statements.
func iciciCreditRow(row []string) (*models.ConvertedTransactionDoc, bool, error) {
	// Due to some reason, the statements contain empty rows in between too.
	if row[0] == "" {
		// Skipping to the next row.
		return nil, false, nil
	}

	// Parse timestamp.
	timestamp, err := time.Parse("02/01/2006", row[0])
	if err != nil {
		// If we fail to parse the timestamp, we consider it as the end of the transaction table.
		return nil, true, nil //nolint:nilerr
	}

	// Other required fields.
	serial, remarks, amountSign := row[1], row[2], row[6]

	// Get the amount information.
	amount, err := strconv.ParseFloat(row[5], 64)
	if err != nil {
		// If we fail to parse the amount, we consider it as the end of the transaction table.
		return nil, true, nil //nolint:nilerr
	}

	// If the amount is zero, we do not consider this row.
	if amount == 0 {
		// Skipping to the next row.
		return nil, false, nil
	}

	// Apparently, debit transactions have "CR" BillingAmountSign -_-
	if amountSign == "CR" {
		amount *= -1
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
