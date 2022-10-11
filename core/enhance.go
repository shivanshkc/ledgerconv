package core

import (
	"context"
)

// Enhance adds the custom fields with zero values to the statement present in the input directory, and places this
// new statement in the output directory.
//
// If there is already a statement in the output directory, then conflicting statement entries are skipped.
//
// This is an idempotent operation.
func Enhance(ctx context.Context, inputDir string, outputDir string) error {
	panic("implement me")
}
