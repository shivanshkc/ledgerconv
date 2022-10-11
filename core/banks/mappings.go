package banks

import (
	"fmt"
	"strings"
)

const (
	iciciSavings BankAccountType = "icici-savings"
	iciciCredit  BankAccountType = "icici-credit"
	hdfcSavings  BankAccountType = "hdfc-savings"
	hdfcCredit   BankAccountType = "hdfc-credit"
	// New account types can be added here.
)

// accountTypeInferRules define how an informal account name will be converted into a definite BankAccountType.
var accountTypeInferRules = msi{
	// If the account name contains the "credit" keyword...
	"credit": msi{
		// And it contains the "icici" keyword as well...
		"icici": iciciCredit, // Then the account type is iciciCredit.
		// And it contains the "hdfc" keyword as well...
		"hdfc": hdfcCredit, // Then the account type is hdfcCredit.
	},
	// If the account name contains the "icici" keyword...
	"icici": iciciSavings, // Then the account type is iciciSavings.
	// If the account name contains the "hdfc" keyword...
	"hdfc": hdfcSavings, // Then the account type is hdfcSavings.
}

// ConverterMap maps bank account types to their respective ConverterFunc.
var ConverterMap = map[BankAccountType]ConverterFunc{
	iciciSavings: convICICISavings,
	iciciCredit:  convICICICredit,
	hdfcSavings:  convHDFCSavings,
	hdfcCredit:   nil, // This means not-implemented.
}

// InferAccountType accepts an account name and infers its type.
func InferAccountType(account string) (BankAccountType, error) {
	// Convert to lower case for case-insensitive comparison.
	account = strings.ToLower(account)
	// Rules are recursive in nature. Hence, we need a recursive function to deal with them.
	return inferAccountTypeRecursive(account, accountTypeInferRules)
}

// inferAccountTypeRecursive loops over the given rules to determine the bank account type.
func inferAccountTypeRecursive(account string, rules msi) (BankAccountType, error) {
	// Loop over the infer-rules to figure out the account type.
	for key, value := range rules {
		// If the keyword does not match, we continue to the next key.
		if !strings.Contains(account, strings.ToLower(key)) {
			continue
		}

		// Check if there are no further nested rules.
		valueBAT, isBAT := value.(BankAccountType)
		if isBAT {
			return valueBAT, nil
		}

		// Check if there are further nested rules.
		valueMSI, isMSI := value.(msi)
		if !isMSI {
			return "", fmt.Errorf("invalid rule structure detected")
		}

		// Recursive call into the nested rules.
		return inferAccountTypeRecursive(account, valueMSI)
	}

	return "", fmt.Errorf("no rules matched")
}
