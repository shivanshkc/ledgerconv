package banks

import (
	"time"
)

// TransactionDoc represents a bank agnostic transaction document.
type TransactionDoc struct {
	AccountName string `json:"account_name"`

	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`

	RefNum      string `json:"bank_ref_num"`
	PaymentMode string `json:"payment_mode"`
	Remarks     string `json:"remarks"`
}

// ConverterFunc represents a function that converts CSV bank statements into a list of transaction documents.
// It can be implemented differently for each bank as per their statement format.
type ConverterFunc func(csvContent [][]string) ([]*TransactionDoc, error)

// BankAccountType is an enum for the type of bank accounts.
type BankAccountType string
