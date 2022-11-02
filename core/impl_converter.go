package core

import (
	"context"
	"fmt"
	"path"
	"sort"

	"github.com/shivanshkc/ledgerconv/core/banks"
	"github.com/shivanshkc/ledgerconv/core/models"
	"github.com/shivanshkc/ledgerconv/core/utils/io"
)

// converter implements the Converter interface.
type converter struct{}

// NewConverter is the constructor for the underlying implementation of the Converter.
func NewConverter() Converter {
	return &converter{}
}

//nolint:funlen,cyclop // Core functions are allowed to be big.
func (c *converter) Convert(ctx context.Context, inputDir string, outputFile string) error {
	// List account directories.
	accountDirs, err := io.ListVisibleDir(inputDir)
	if err != nil {
		return fmt.Errorf("failed to list account directories: %w", err)
	}

	// Convert the account directory names to account types.
	accountTypes := make([]models.BankAccountType, len(accountDirs))
	for idx, dir := range accountDirs {
		// Business call to infer the account type.
		accountType, err := banks.InferAccountType(dir, nil)
		if err != nil {
			return fmt.Errorf("failed to infer account type for: %s, because: %w", dir, err)
		}
		// Collect the inferred type into the declared slice.
		accountTypes[idx] = accountType
	}

	// This will collect all the converted transaction documents.
	var statement []*models.ConvertedTransactionDoc

	// Loop over all account directories to read their CSV files and convert them.
	for idx, dir := range accountDirs {
		// Form the full path to the account directory.
		accountDirPath := path.Join(inputDir, dir)
		// List CSV files in the account directory.
		csvFiles, err := io.ListCSVFiles(accountDirPath)
		if err != nil {
			return fmt.Errorf("failed to list csv files in %s, because: %w", accountDirPath, err)
		}

		// Get the converter function for this account type.
		converterFunc, exists := banks.AccountTypeConverterMap[accountTypes[idx]]
		if !exists || converterFunc == nil {
			return fmt.Errorf("bank account: %s is not supported", dir)
		}

		// Loop over the CSV files of this account directory to convert them.
		for _, csv := range csvFiles {
			// Form the full path to the csv file.
			csvPath := path.Join(accountDirPath, csv)
			// Get the contents of the csv file.
			csvContent, err := io.ReadWholeCSV(csvPath)
			if err != nil {
				return fmt.Errorf("failed to read csv file: %s, because: %w", csvPath, err)
			}

			// Convert the csv content to list of transactions.
			transactions, err := converterFunc(csvContent)
			if err != nil {
				return fmt.Errorf("failed to convert: %s to statement, because %w", csvPath, err)
			}

			// Add account name to all transactions.
			for _, txn := range transactions {
				txn.AccountName = dir
			}

			// Collect transactions in the main slice.
			statement = append(statement, transactions...)
		}
	}

	// Sort the converted statement.
	sort.SliceStable(statement, func(i, j int) bool {
		return statement[i].Timestamp.After(statement[j].Timestamp)
	})

	// Write statement file.
	if err := io.WriteJSONFile(outputFile, statement); err != nil {
		return fmt.Errorf("failed to write  converted statement file: %w", err)
	}

	return nil
}
