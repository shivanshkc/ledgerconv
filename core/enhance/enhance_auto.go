package enhance

import (
	"strings"

	"github.com/shivanshkc/ledgerconv/core/models"
	"github.com/shivanshkc/ledgerconv/core/utils"
)

// Auto attempts to auto-enhance the given converted transaction with the help of the given auto-enhance spec.
func Auto(txn *models.ConvertedTransactionDoc, spec []*models.AutoEnhanceSpec) (
	*models.EnhancedTransactionDoc, bool, error,
) {
	// -------------------------------------------
	// Quantity to be ultimately returned.
	enhanced := &models.EnhancedTransactionDoc{
		ConvertedTransactionDoc: txn,
		Categories:              &models.AmountPerCategory{},
		Labels:                  nil,
		Summary:                 "",
	}

	// Convert to lower case for case-insensitive comparison.
	bankRemarks := strings.ToLower(txn.BankRemarks)

	// Check against each element.
	for _, elem := range spec {
		// Respect the ForCredit flag.
		if (elem.ForCredit && txn.Amount < 0) || (!elem.ForCredit && txn.Amount > 0) {
			continue
		}

		// If no matches, we move to the next iteration.
		if !utils.ContainsAnyNoCase(bankRemarks, elem.RemarksKeywords) {
			continue
		}

		// Convert categories percentage to actual values.
		// TODO: This is way too manual. Find out a better way.
		if txn.Amount > 0 {
			enhanced.Categories.Salary = elem.Categories.Salary * txn.Amount / 100.0
			enhanced.Categories.Returns = elem.Categories.Returns * txn.Amount / 100.0
			enhanced.Categories.Misc = elem.Categories.Misc * txn.Amount / 100.0
		} else {
			enhanced.Categories.Essentials = elem.Categories.Essentials * txn.Amount / 100.0
			enhanced.Categories.Investments = elem.Categories.Investments * txn.Amount / 100.0
			enhanced.Categories.Savings = elem.Categories.Savings * txn.Amount / 100.0
			enhanced.Categories.Luxury = elem.Categories.Luxury * txn.Amount / 100.0
		}
		enhanced.Categories.Ignorable = elem.Categories.Ignorable * txn.Amount / 100.0

		// Other fields.
		enhanced.Labels = elem.Labels
		enhanced.Summary = elem.Summary
		enhanced.AutoEnhanced = true

		return enhanced, true, nil
	}

	return nil, false, nil
}
