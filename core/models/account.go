package models

// BankAccountType is an enum for the type of bank accounts.
type BankAccountType string

const (
	// ICICISavings is the BankAccountType for an ICICI savings account.
	ICICISavings BankAccountType = "icici-savings"
	// ICICICredit is the BankAccountType for an ICICI credit card account.
	ICICICredit BankAccountType = "icici-credit"
	// HDFCSavings is the BankAccountType for an HDFC savings account.
	HDFCSavings BankAccountType = "hdfc-savings"
	// HDFCCredit is the BankAccountType for an HDFC credit card account.
	HDFCCredit BankAccountType = "hdfc-credit"
)
