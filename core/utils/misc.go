package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
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

// ContainsAnyNoCase checks if the provided mainStr contains any of the subStrings, case-insensitively.
func ContainsAnyNoCase(mainStr string, subStrings []string) bool {
	// Convert to lower case for case-insensitive matching.
	mainStr = strings.ToLower(mainStr)

	// Loop over all provided sub-strings to find matches.
	for _, sub := range subStrings {
		// Convert to lower case for case-insensitive matching.
		sub = strings.ToLower(sub)
		if strings.Contains(mainStr, sub) {
			return true
		}
	}

	return false
}
