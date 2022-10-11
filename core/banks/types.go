package banks

import (
	"github.com/shivanshkc/ledgerconv/core/models"
)

// Shorthand for map[string]interface{}.
type msi = map[string]interface{}

// ConverterFunc represents a function that converts CSV bank statements into a list of transaction documents.
// It can be implemented differently for each bank as per their statement format.
type ConverterFunc func(csvContent [][]string) ([]*models.TransactionDoc, error)

// BankAccountType is an enum for the type of bank accounts.
type BankAccountType string
