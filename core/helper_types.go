package core

import (
	"time"
)

// transactionDoc represents a bank agnostic transaction document.
type transactionDoc struct {
	AccountName string `json:"account_name"`

	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`

	RefNum      string `json:"bank_ref_num"`
	PaymentMode string `json:"payment_mode"`
	Remarks     string `json:"remarks"`
}

// converterFunc represents a function that converts CSV bank statements into a list of transaction documents.
// It can be implemented differently for each bank as per their statement format.
type converterFunc func(csvContent [][]string) ([]*transactionDoc, error)

// bankAccountType is an enum for the type of bank accounts.
type bankAccountType string

const (
	iciciSavings bankAccountType = "icici-savings"
	iciciCredit  bankAccountType = "icici-credit"

	hdfcSavings bankAccountType = "hdfc-savings"
	hdfcCredit  bankAccountType = "hdfc-credit"
)
