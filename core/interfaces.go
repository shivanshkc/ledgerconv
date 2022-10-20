package core

import (
	"context"
)

// Converter converts CSV bank statements (also called "Original" statements) into a single bank-agnostic JSON
// statement (also called "Converted" statement).
//
// Inputs:
//
//  1. inputDir - where the account directories containing the original statements are present.
//  2. outputFile - where the converted statement file will be stored.
type Converter interface {
	Convert(ctx context.Context, inputDir string, outputFile string) error
}

// Enhancer converts a converted statement into an enhanced statement.
//
// Enhanced statements are nothing but converted statements with some additional fields that are useful for statistics.
//
// Inputs:
//
//  1. inputFile - where the converted statement file is present.
//  2. outputFile - where the enhanced statement file is present or will be stored.
//  3. autoEnhanceSpecFile - Path to the specification file that allows auto-enhancement.
type Enhancer interface {
	Enhance(ctx context.Context, inputFile string, outputFile string, autoEnhanceSpecFile string) error
}
