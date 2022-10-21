package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// Checksum provides the checksum of the given object.
func Checksum(input interface{}) (string, error) {
	// Convert to byte slice to write to hash.
	marshalled, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal object because: %w", err)
	}
	// Calculate, format and return the checksum.
	return fmt.Sprintf("%x", sha256.Sum256(marshalled)), nil
}
