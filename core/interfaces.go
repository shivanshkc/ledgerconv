package core

import (
	"context"
)

// Converter converts CSV bank statements (also called "Original" statements) into a single bank-agnostic JSON
// statement (also called "Converted" statement).
//
// Inputs:
//
//  1. inputPath - where the account directories containing the Original statements are present.
//  2. outputPath - where the Converted statement will be stored.
type Converter interface {
	Convert(ctx context.Context, inputPath string, outputPath string) error
}

// Enhancer converts a Converted statement into an Enhanced statement.
//
// Enhanced statements are nothing but Converted statements with some additional fields that are useful for statistics.
//
// Inputs:
//
//  1. inputPath - where the Converted statement is present.
//  2. outputPath - where the Enhanced statement is present or will be stored.
//  3. autoEnhanceSpecPath - Path to the specification that allows auto-enhancement.
type Enhancer interface {
	Enhance(ctx context.Context, inputPath string, outputPath string, autoEnhanceSpecPath string) error
}
