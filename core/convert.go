package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

// converterMap maps bank account types to their respective converterFuncs.
var converterMap = map[bankAccountType]converterFunc{
	iciciSavings: convICICISavings,
	iciciCredit:  convICICICredit,
	hdfcSavings:  convHDFCSavings,
	hdfcCredit:   convHDFCCredit,
}

// Convert converts all the bank statements in the inputDir into JSON format and stores them into the outputDir.
//
//nolint:funlen,cyclop // Core functions are allowed to be big.
func Convert(ctx context.Context, inputDir string, outputDir string) error {
	// List all account directories.
	accountDirs, err := showDirs(inputDir)
	if err != nil {
		return fmt.Errorf("failed to list accounts directories: %w", err)
	}

	// All transactions will be collected in this slice.
	var transactionDocs []*transactionDoc

	// Loop over all account directories to convert all their statements.
	for _, accountDir := range accountDirs {
		// Infer the account type for this account. This is needed to pick the right converterFunc.
		accountType, err := inferAccountType(accountDir)
		if err != nil {
			return fmt.Errorf("failed to infer account type for account: %s, because: %w", accountDir, err)
		}

		// Pick the right converterFunc for this account.
		converter, exists := converterMap[accountType]
		if !exists || converter == nil {
			return fmt.Errorf("no converterFunc implementation found for this account type: %s, for directory: %s",
				accountType, accountDir)
		}

		// Complete path to this account directory.
		statementDir := path.Join(inputDir, accountDir)
		// List all the statement files.
		statementFiles, err := showFiles(statementDir)
		if err != nil {
			return fmt.Errorf("failed to list statement files in directory: %s, because: %w", statementDir, err)
		}

		// Loop over each statement file to convert it.
		for _, statementFile := range statementFiles {
			// Complete path to the statement file.
			pathToFile := path.Join(statementDir, statementFile)
			// Read the statement file for conversion.
			csvContent, err := readCSV(pathToFile)
			if err != nil {
				return fmt.Errorf("failed to read the statement file: %s, because: %w", pathToFile, err)
			}

			// Convert the CSV content into transaction list.
			txDocs, err := converter(csvContent)
			if err != nil {
				return fmt.Errorf("failed to convert statement file: %s, because: %w", pathToFile, err)
			}

			// Add account name to all transactions.
			for _, doc := range txDocs {
				doc.AccountName = accountDir
			}

			// Collect results.
			transactionDocs = append(transactionDocs, txDocs...)
		}
	}

	// Marshal the transaction list to write into file.
	transactionDocsBytes, err := json.Marshal(transactionDocs)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction list: %w", err)
	}

	// Name of the output file.
	outputFilePath := path.Join(outputDir, fmt.Sprintf("transactions-%d.json", time.Now().Unix()))
	// Write the output file.
	if err := os.WriteFile(outputFilePath, transactionDocsBytes, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// inferAccountType accepts an account name and infers its type.
func inferAccountType(account string) (bankAccountType, error) {
	// These keywords can be modified to control the behaviour of this function.
	creditCardKeywords := []string{"credit"}
	iciciKeywords := []string{"icici"}
	hdfcKeywords := []string{"hdfc"}

	// Check if the account is a credit card account.
	isCreditCard := containsAnyCaseInsensitive(account, creditCardKeywords)

	if containsAnyCaseInsensitive(account, iciciKeywords) {
		if isCreditCard {
			return iciciCredit, nil
		}
		return iciciSavings, nil
	}

	if containsAnyCaseInsensitive(account, hdfcKeywords) {
		if isCreditCard {
			return hdfcCredit, nil
		}
		return hdfcSavings, nil
	}

	return "", fmt.Errorf("no keyword matches found")
}
