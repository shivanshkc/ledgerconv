package models

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
