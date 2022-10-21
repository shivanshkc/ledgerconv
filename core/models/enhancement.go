package models

// AutoEnhanceSpec is the schema of an element in an auto-enhance-spec file.
type AutoEnhanceSpec struct {
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
