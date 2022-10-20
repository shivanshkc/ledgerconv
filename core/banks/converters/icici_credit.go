package converters

import (
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
	return nil, false, nil
}
