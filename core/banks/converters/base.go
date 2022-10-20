package converters

import (
	"github.com/shivanshkc/ledgerconv/core/models"
)

// rowConvFunc represents a function that converts a single CSV row to a converted transaction document.
// It is meant to be implemented separately for different banks.
type rowConvFunc func(row []string) (*models.ConvertedTransactionDoc, error)

// base is the base converter that all other converter function can rely upon.
func base(header []string, startingOffset int, convFunc rowConvFunc) ([]*models.ConvertedTransactionDoc, error) {
	return nil, nil
}
