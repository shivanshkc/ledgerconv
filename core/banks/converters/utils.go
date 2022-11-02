package converters

import (
	"reflect"
	"strings"
)

// trimCSV trims all elements of the given csvContent.
func trimCSV(csvContent [][]string) [][]string {
	for i := range csvContent {
		for j := range csvContent[i] {
			csvContent[i][j] = strings.TrimSpace(csvContent[i][j])
		}
	}
	return csvContent
}

// getHeaderIndex provides the index of the given header in the given csvContent.
func getHeaderIndex(header []string, csvContent [][]string) int {
	// Loop over CSV rows to find the starting of the transaction table.
	for idx, row := range csvContent {
		// If the row is not the same as the starting header, we continue.
		if !reflect.DeepEqual(row, header) {
			continue
		}
		// Starting of the transaction table is located.
		return idx
	}

	return 0
}
