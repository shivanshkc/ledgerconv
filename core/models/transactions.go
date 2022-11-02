package models

import (
	"time"
)

// ConvertedTransactionDoc represents a bank agnostic transaction document.
type ConvertedTransactionDoc struct {
	// AccountName is the name of the account to which the transaction belongs.
	AccountName string `json:"account_name"`

	// Amount of the transaction.
	Amount float64 `json:"amount"`
	// Timestamp of the transaction.
	Timestamp time.Time `json:"timestamp"`

	// BankSerial can be the transaction ID or transaction reference number, as mentioned in the bank statement.
	BankSerial string `json:"bank_ref_num"`
	// BankPaymentMode is the payment mode of this transaction, as mentioned in the bank statement.
	BankPaymentMode string `json:"bank_payment_mode"`
	// BankRemarks are the notes/narration/remarks of this transaction, as mentioned in the bank statement.
	BankRemarks string `json:"bank_remarks"`
}

// EnhancedTransactionDoc is a super set of the ConvertedTransactionDoc. It contains extra fields to persist more
// useful information about the transaction.
type EnhancedTransactionDoc struct {
	*ConvertedTransactionDoc

	// DocCorrelationID correlates this doc with its corresponding ConvertedTransactionDoc.
	DocCorrelationID string `json:"doc_correlation_id"`

	// Categories tells how the amount is distributed amongst the different categories.
	Categories *AmountPerCategory `json:"categories"`
	// Labels of the transaction.
	Labels []string `json:"labels"`
	// Summary of the transaction.
	Summary string `json:"summary"`

	// AutoEnhanced is a flag that tells whether the transaction was auto-enhanced or not.
	AutoEnhanced bool `json:"auto_enhanced"`
}

// AmountPerCategory holds the distribution of a transaction's amount over all categories.
type AmountPerCategory struct {
	// DEBIT CATEGORIES ##############################

	// Essentials are those debits that a person cannot avoid. Example: House EMI, electricity bills, anniversaries.
	Essentials float64 `json:"essentials"`
	// Investments can be stocks, equity, real-estate, crypto etc.
	Investments float64 `json:"investments"`
	// Luxury is money that is deliberately spent on comforts.
	Luxury float64 `json:"luxury"`
	// Savings are required in case of an immediate emergency.
	Savings float64 `json:"savings"`
	// ###############################################

	// CREDIT CATEGORIES #############################

	// Salary is primary source of income.
	Salary float64 `json:"salary"`
	// Returns can be any investment return, including bank account interest.
	Returns float64 `json:"returns"`
	// Misc are all other kinds of income. Including petty credits.
	Misc float64 `json:"misc"`
	// ###############################################

	// COMMON CATEGORIES #############################

	// Ignorable contains those transactions that add up to zero, and hence should not contribute to any stat.
	Ignorable float64 `json:"ignorable"`
	// ###############################################
}

// HasOnlyDebit returns true if the AmountPerCategory only has debit categories filled.
func (a *AmountPerCategory) HasOnlyDebit() bool {
	return a.Salary == 0 && a.Returns == 0 && a.Misc == 0
}

// HasOnlyCredit returns true if the AmountPerCategory only has credit categories filled.
func (a *AmountPerCategory) HasOnlyCredit() bool {
	return a.Essentials == 0 && a.Investments == 0 && a.Luxury == 0 && a.Savings == 0
}

// DebitSum returns the sum of all debit category values.
func (a *AmountPerCategory) DebitSum() float64 {
	return a.Essentials + a.Investments + a.Luxury + a.Savings + a.Ignorable
}

// CreditSum returns the sum of all credit category values.
func (a *AmountPerCategory) CreditSum() float64 {
	return a.Salary + a.Returns + a.Misc + a.Ignorable
}
