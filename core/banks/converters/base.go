package converters

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/shivanshkc/ledgerconv/core/models"
)

// rowConvFunc represents a function that converts a single CSV row to a converted transaction document.
// It is meant to be implemented separately for different banks.
type rowConvFunc func(row []string) (converted *models.ConvertedTransactionDoc, done bool, err error)

// base is the base converter that all other converter function can rely upon.
//
//nolint:cyclop // TODO?
func base(content [][]string, header []string, headerRowOffset int, convFunc rowConvFunc) (
	[]*models.ConvertedTransactionDoc, error,
) {
	// ----------------------------------------
	// Trim each element. Bank statement schemas are not to be trusted!
	for i := range content {
		for j := range content[i] {
			content[i][j] = strings.TrimSpace(content[i][j])
		}
	}

	// This var will hold the index of the first transaction table row.
	var startingIdx int

	// Loop over CSV rows to find the starting of the transaction table.
	for idx, row := range content {
		// If the row is not the same as the starting header, we continue.
		if !reflect.DeepEqual(row, header) {
			continue
		}
		// Starting of the transaction table is located.
		startingIdx = idx + headerRowOffset
		break
	}

	// Just a safety check.
	if startingIdx == 0 || startingIdx >= len(content) {
		return nil, nil
	}

	// This var will hold the final list of converted transactions.
	var statement []*models.ConvertedTransactionDoc

	// Begin looping over the transaction table.
	for _, row := range content[startingIdx:] {
		// Convert the row.
		converted, done, err := convFunc(row)
		if err != nil {
			return nil, fmt.Errorf("failed to parse statement row: %w", err)
		}
		// If the converted document is non-nil, we add it to the main slice.
		if converted != nil {
			statement = append(statement, converted)
		}
		// Check if the table is exhausted.
		if done {
			break
		}
	}

	return statement, nil
}
