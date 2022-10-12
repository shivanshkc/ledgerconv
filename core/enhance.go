package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/shivanshkc/ledgerconv/core/models"
)

// enhancedFilename is the name of the file in which the enhanced transactions will be written.
const enhancedFilename = "enhanced-transactions.json"

// Enhance adds the custom fields with zero values to the statement present in the input directory, and places this
// new statement in the output directory.
//
// If there is already a statement in the output directory, then conflicting statement entries are skipped.
//
// This is an idempotent operation.
//
//nolint:cyclop // Core functions are allowed to be big.
func Enhance(ctx context.Context, inputDir string, outputDir string) error {
	// Path to the existing enhanced statement file.
	enhancedFilepath := path.Join(outputDir, enhancedFilename)
	// Open the enhanced statement file to load existing enhanced transactions.
	enhancedFileReader, err := os.Open(enhancedFilepath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to open file: %s, because: %w", enhancedFilepath, err)
	}

	// Decode the enhanced statement into a slice. If the file did not exist, the reader will be nil, and no decoding
	// will take place.
	var enhancedStatement []*models.EnhancedTransactionDoc
	if enhancedFileReader != nil {
		defer func() { _ = enhancedFileReader.Close() }()

		if err := json.NewDecoder(enhancedFileReader).Decode(&enhancedStatement); err != nil {
			return fmt.Errorf("failed to decode the enhanced statement: %w", err)
		}
	}

	// This maps the enhanced transactions to their correlation IDs.
	enhancedStatementMap := map[string]*models.EnhancedTransactionDoc{}
	for _, tx := range enhancedStatement {
		enhancedStatementMap[tx.DocCorrelationID] = tx
	}

	// Path to the converted statement file.
	convertedFilepath := path.Join(inputDir, convertedFilename)
	// Open the converted statement file to enhance the transactions.
	convertedFileReader, err := os.Open(convertedFilepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %s, because: %w", convertedFilepath, err)
	}
	defer func() { _ = convertedFileReader.Close() }()

	// Decode the converted statement into a slice.
	var convertedStatement []*models.ConvertedTransactionDoc
	if err := json.NewDecoder(convertedFileReader).Decode(&convertedStatement); err != nil {
		return fmt.Errorf("failed to decode the converted statement: %w", err)
	}

	// Loop over all converted transactions to enhance them.
	for _, txn := range convertedStatement {
		// Generate checksum.
		correlationID, err := genConvertedTxChecksum(txn)
		if err != nil {
			return fmt.Errorf("failed to generate checksum: %w", err)
		}
		// See if enhanced transaction exists.
		if _, exists := enhancedStatementMap[correlationID]; exists {
			continue
		}

		fmt.Println(">> Creating transaction with correlation ID:", correlationID)
		// Prompt the user for inputs.
		_, _, _ = takeUserInput(txn)
	}

	return nil
}

// takeUserInput prompts the user for inputs required to create an enhanced transaction.
func takeUserInput(tx *models.ConvertedTransactionDoc) (*models.AmountPerCategory, []string, string) {
	return nil, nil, ""
}
