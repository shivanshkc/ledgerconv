package banks

import (
	"fmt"
	"strings"

	"github.com/shivanshkc/ledgerconv/core/banks/converters"
	"github.com/shivanshkc/ledgerconv/core/models"
)

// msi is a simple type-alias for map[string]interface{}.
type msi = map[string]interface{}

// ConverterFunc represents a function that converts CSV bank statements into a list of transaction documents.
// It can be implemented differently for each bank as per their statement format.
type ConverterFunc func(csvContent [][]string) ([]*models.ConvertedTransactionDoc, error)

// accountTypeInferRules define how an informal account name will be converted into a definite account type.
var accountTypeInferRules = msi{
	// If the account name contains the "credit" keyword...
	"credit": msi{
		// And it contains the "icici" keyword as well...
		"icici": models.ICICICredit, // Then the account type is iciciCredit.
		// And it contains the "hdfc" keyword as well...
		"hdfc": models.HDFCCredit, // Then the account type is hdfcCredit.
	},
	// If the account name contains the "icici" keyword...
	"icici": models.ICICISavings, // Then the account type is iciciSavings.
	// If the account name contains the "hdfc" keyword...
	"hdfc": models.HDFCSavings, // Then the account type is hdfcSavings.
}

// AccountTypeConverterMap maps bank account types to their respective ConverterFunc.
var AccountTypeConverterMap = map[models.BankAccountType]ConverterFunc{
	models.ICICISavings: converters.ICICISavings,
	models.ICICICredit:  converters.ICICICredit,
	models.HDFCSavings:  converters.HDFCSavings,
	models.HDFCCredit:   nil, // Not implemented yet.
}

// InferAccountType accepts an account name and infers its account type.
func InferAccountType(account string, rules msi) (models.BankAccountType, error) {
	// Convert to lower case for case-insensitive comparison.
	account = strings.ToLower(account)
	// If nil rules are provided, it means it is not a recursive call.
	if rules == nil {
		rules = accountTypeInferRules
	}

	// Loop over the infer-rules to figure out the account type.
	for key, value := range rules {
		// If the keyword does not match, we continue to the next key.
		if !strings.Contains(account, strings.ToLower(key)) {
			continue
		}

		// Check if there are no further nested rules.
		valueBAT, isBAT := value.(models.BankAccountType)
		if isBAT {
			return valueBAT, nil
		}

		// Check if there are further nested rules.
		valueMSI, isMSI := value.(msi)
		if !isMSI {
			return "", fmt.Errorf("invalid rule structure detected")
		}

		// Recursive call into the nested rules.
		return InferAccountType(account, valueMSI)
	}

	// Failed to infer account type.
	return "", fmt.Errorf("no rules matched")
}
