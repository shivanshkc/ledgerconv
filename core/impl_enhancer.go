package core

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/shivanshkc/ledgerconv/core/models"
	"github.com/shivanshkc/ledgerconv/core/utils"
	"github.com/shivanshkc/ledgerconv/core/utils/io"

	"github.com/fatih/color"
)

// enhancer implements the Enhancer interface.
type enhancer struct{}

// NewEnhancer is the constructor for the underlying implementation of the Enhancer.
func NewEnhancer() Enhancer {
	return &enhancer{}
}

//nolint:cyclop // Core functions are allowed to be big.
func (e *enhancer) Enhance(ctx context.Context, inputFile string, outputFile string, specFile string) error {
	// Read the provided converted statement.
	var convertedStm []*models.ConvertedTransactionDoc
	if err := io.ReadJSONFile(inputFile, &convertedStm); err != nil {
		return fmt.Errorf("failed to read the converted statement at: %s, because: %w", inputFile, err)
	}

	// If the converted statement is empty, we don't have anything to do.
	if len(convertedStm) == 0 {
		return fmt.Errorf("converted statement is empty")
	}

	// Read the provided enhanced statement. If the file does not exist, it will be ignored.
	var enhancedStm []*models.EnhancedTransactionDoc
	if err := io.ReadJSONFile(outputFile, &enhancedStm); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to read the enhanced statement at: %s, because: %w", outputFile, err)
	}

	// This maps the enhanced transactions to their correlation IDs.
	enhancedStmMap := map[string]*models.EnhancedTransactionDoc{}
	for _, tx := range enhancedStm {
		enhancedStmMap[tx.DocCorrelationID] = tx
	}

	// Read the provided auto-enhance spec file.
	var autoEnhSpec []*models.AutoEnhanceSpec
	if err := io.ReadJSONFile(specFile, &autoEnhSpec); err != nil {
		return fmt.Errorf("failed to read the auto-enhancement spec file at: %s, because: %w", specFile, err)
	}

	// This will hold only those converted transactions that have not already been enhanced.
	var newConvertedStm []*models.ConvertedTransactionDoc
	// Loop over all converted transaction to filter out the already enhanced ones.
	for _, txn := range convertedStm {
		// Get transaction checksum.
		checksum, err := utils.Checksum(txn)
		if err != nil {
			return fmt.Errorf("failed to calculate checksum for: %+v, because: %w", txn, err)
		}
		// See if enhanced transaction exists.
		if _, exists := enhancedStmMap[checksum]; !exists {
			newConvertedStm = append(newConvertedStm, txn)
		}
	}

	// The enhancement loop.
	for idx, txn := range newConvertedStm {
		color.Yellow("=================================================================")
		color.Yellow(fmt.Sprintf("Processing transaction %d out of %d", idx+1, len(newConvertedStm)))

		_ = txn
	}

	return nil
}
