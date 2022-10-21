package models

import (
	"fmt"
)

// AutoEnhanceSpec is the schema of an element in an auto-enhance-spec file.
type AutoEnhanceSpec struct {
	// ForCredit is a flag. If it's set to true, this spec element will only be applied to credit transactions.
	// Otherwise, to only debit transactions.
	ForCredit bool `json:"for_credit"`
	// RemarksKeywords are the keywords that the bank-remarks of the transaction should match for auto-enhancement.
	RemarksKeywords []string `json:"remarks_keywords"`
	// Categories is the intended amount-category distribution for the transaction.
	// The values should be in percentages.
	Categories *AmountPerCategory `json:"categories"`
	// Labels are the intended labels of the transaction.
	Labels []string `json:"labels"`
	// Summary is the intended summary of the transaction.
	Summary string `json:"summary"`
}

// Validate returns any validation errors in the AutoEnhanceSpec.
func (a *AutoEnhanceSpec) Validate() error {
	if a.ForCredit && !(a.Categories.HasOnlyCredit() && a.Categories.CreditSum() == 100) {
		return fmt.Errorf("if 'for_credit' is true, then 'categories' should only contain credit categories " +
			" and their sum should equal 100")
	}

	if !a.ForCredit && !(a.Categories.HasOnlyDebit() && a.Categories.DebitSum() == 100) {
		return fmt.Errorf("if 'for_credit' is false, then 'categories' should only contain debit categories " +
			" and their sum should equal 100")
	}

	return nil
}
