package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/shivanshkc/ledgerconv/core/enhance"
	"github.com/shivanshkc/ledgerconv/core/models"
	"github.com/shivanshkc/ledgerconv/core/utils"
	"github.com/shivanshkc/ledgerconv/core/utils/io"

	"github.com/fatih/color"
)

const defaultSpecFile = "./auto-enhance-spec.json"

// enhancer implements the Enhancer interface.
type enhancer struct{}

// NewEnhancer is the constructor for the underlying implementation of the Enhancer.
func NewEnhancer() Enhancer {
	return &enhancer{}
}

//nolint:funlen,cyclop // Core functions are allowed to be big.
func (e *enhancer) Enhance(ctx context.Context, inputFile string, outputFile string, specFile string, onlyAuto bool,
) error {
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
	var enhancedStm []*models.EnhancedTransactionDoc //nolint:prealloc // Cannot preallocate.
	if err := io.ReadJSONFile(outputFile, &enhancedStm); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to read the enhanced statement at: %s, because: %w", outputFile, err)
	}

	// This maps the enhanced transactions to their correlation IDs.
	enhancedStmMap := map[string]*models.EnhancedTransactionDoc{}
	for _, tx := range enhancedStm {
		enhancedStmMap[tx.DocCorrelationID] = tx
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

	// Resolve the auto-enhance spec file.
	autoEnhSpec, err := resolveSpecFile(specFile)
	if err != nil {
		return fmt.Errorf("failed to resolve the auto-enhance spec file: %w", err)
	}

	// The enhancement loop.
	for idx, txn := range newConvertedStm {
		color.Cyan("=================================================================")
		color.Blue(fmt.Sprintf("Processing transaction %d out of %d", idx+1, len(newConvertedStm)))

		// Attempt to auto-enhance.
		enhancedTx, done, err := enhance.Auto(txn, autoEnhSpec)
		if err != nil {
			return fmt.Errorf("failed to auto-enhance transaction: %+v, because: %w", txn, err)
		}

		//nolint:gocritic // Cannot write as switch statement.
		if done {
			color.Cyan("-----------------------------------------------------------------")
			color.Green("Auto enhanced.")
		} else if !onlyAuto {
			// Enhance manually.
			enhancedTx, err = enhance.Manual(txn)
			if err != nil {
				return fmt.Errorf("failed to enhance transaction: %+v, because: %w", txn, err)
			}
		} else {
			color.Cyan("-----------------------------------------------------------------")
			color.Blue("Not auto-enhanceable. Skipped.")
			continue
		}

		// Generate checksum.
		checksum, err := utils.Checksum(txn)
		if err != nil {
			return fmt.Errorf("failed to generate checksum for tx: %+v, because: %w", txn, err)
		}

		// Make sure these fields remain the same before and after enhancement.
		enhancedTx.ConvertedTransactionDoc = txn
		enhancedTx.DocCorrelationID = checksum

		// Update the main statement.
		enhancedStm = append(enhancedStm, enhancedTx)

		// Sort the enhanced statement.
		sort.SliceStable(enhancedStm, func(i, j int) bool {
			return enhancedStm[i].Timestamp.After(enhancedStm[j].Timestamp)
		})

		// Write statement file.
		if err := io.WriteJSONFile(outputFile, enhancedStm); err != nil {
			return fmt.Errorf("failed to write enhanced statement file: %w", err)
		}

		color.Cyan("-----------------------------------------------------------------")
		color.Green("Saved.")
	}

	return nil
}

// resolveSpecFile is responsible to treat the auto-enhance spec file as an optional parameter.
//
// If the user has provided a path, then it must exist and should contain a valid spec file.
// If the user has not provided a path, then we use try to use the default one, only if it exists.
func resolveSpecFile(filePath string) ([]*models.AutoEnhanceSpec, error) {
	// This keeps track of whether we are using the user given file or the default one.
	var usingDefault bool

	// If the user didn't provide a path, we use the default one and mark the flag.
	if filePath == "" {
		usingDefault, filePath = true, defaultSpecFile
	}

	// Read the spec file.
	var autoEnhSpec []*models.AutoEnhanceSpec
	if err := io.ReadJSONFile(filePath, &autoEnhSpec); err != nil {
		// If we are using the default file, and it does not exist, we ignore the error.
		if errors.Is(err, os.ErrNotExist) && usingDefault {
			return nil, nil
		}
		// Failed to read the file.
		return nil, fmt.Errorf("failed to read the spec file at: %s, because: %w", filePath, err)
	}

	// Check if all elements in the spec file are valid.
	for idx, elem := range autoEnhSpec {
		if err := elem.Validate(); err != nil {
			return nil, fmt.Errorf("invalid auto-enhance spec file: element no: %d, because: %w", idx, err)
		}
	}

	color.Blue("Using auto-enhance spec file: %s", filePath)
	return autoEnhSpec, nil
}
