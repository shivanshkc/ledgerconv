package core

import (
	"context"
)

// converter implements the Converter interface.
type converter struct{}

// NewConverter is the constructor for the underlying implementation of the Converter.
func NewConverter() Converter {
	return &converter{}
}

func (c *converter) Convert(ctx context.Context, inputPath string, outputPath string) error {
	return nil
}
