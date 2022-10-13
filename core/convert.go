package core

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/shivanshkc/ledgerconv/core/banks"
	"github.com/shivanshkc/ledgerconv/core/models"
)

// convertedFilename is the name of the file in which the converted transactions will be written.
const convertedFilename = "converted-transactions.json"

// Convert converts all the bank statements in the inputDir into JSON format and stores them into the outputDir.
//
// This is an idempotent operation.
//
//nolint:funlen,cyclop // Core functions are allowed to be big.
func Convert(ctx context.Context, inputDir string, outputDir string) error {
	// List all account directories.
	accountDirs, err := showDirs(inputDir)
	if err != nil {
		return fmt.Errorf("failed to list accounts directories: %w", err)
	}

	// All transactions will be collected in this slice.
	var transactionDocs []*models.ConvertedTransactionDoc

	// Loop over all account directories to convert all their statements.
	for _, accountDir := range accountDirs {
		// Hidden directories are ignored.
		if strings.HasPrefix(accountDir, ".") {
			continue
		}

		// Infer the account type for this account. This is needed to pick the right converterFunc.
		accountType, err := banks.InferAccountType(accountDir)
		if err != nil {
			return fmt.Errorf("failed to infer account type for account: %s, because: %w", accountDir, err)
		}

		// Pick the right converterFunc for this account.
		converter, exists := banks.ConverterMap[accountType]
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

	// Writing statement file.
	if err := writeJSON(transactionDocs, path.Join(outputDir, convertedFilename)); err != nil {
		return fmt.Errorf("failed to write statement file: %w", err)
	}

	return nil
}
