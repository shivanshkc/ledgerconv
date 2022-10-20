package core

import (
	"context"
)

// enhancer implements the Enhancer interface.
type enhancer struct{}

// NewEnhancer is the constructor for the underlying implementation of the Enhancer.
func NewEnhancer() Enhancer {
	return &enhancer{}
}

func (e *enhancer) Enhance(ctx context.Context, inputPath string, outputPath string, specPath string) error {
	return nil
}
