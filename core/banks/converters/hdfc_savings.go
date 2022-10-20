package converters

import (
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
	return nil, false, nil
}
