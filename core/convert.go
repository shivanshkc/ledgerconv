package core

import (
	"context"
	"fmt"
)

// Convert converts all the bank statements in the inputDir into JSON format and stores them into the outputDir.
func Convert(ctx context.Context, inputDir string, outputDir string) error {
	fmt.Printf("Working with: input: %s, output: %s\n", inputDir, outputDir)
	return nil
}
