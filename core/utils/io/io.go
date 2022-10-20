package io

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// ListVisibleDir lists all visible directories in the provided path.
func ListVisibleDir(dirPath string) ([]string, error) {
	// Read the given directory to get its children.
	elements, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %s, because: %w", dirPath, err)
	}

	// This slice will hold the directories to be returned.
	var visibleDir []string
	// Filter elements.
	for _, elem := range elements {
		// If the element is a directory, and it is visible, it is considered.
		if elem.IsDir() && !strings.HasPrefix(elem.Name(), ".") {
			visibleDir = append(visibleDir, elem.Name())
		}
	}

	return visibleDir, nil
}

// ListCSVFiles lists all CSV files in the provided path.
func ListCSVFiles(dirPath string) ([]string, error) {
	// Read the given directory to get its children.
	elements, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %s, because: %w", dirPath, err)
	}

	// This slice will hold the files to be returned.
	var csvFiles []string
	// Filter elements.
	for _, elem := range elements {
		// If the element is not a directory, and its extension is .csv, it is considered.
		if !elem.IsDir() && strings.HasSuffix(elem.Name(), ".csv") {
			csvFiles = append(csvFiles, elem.Name())
		}
	}

	return csvFiles, nil
}

// ReadWholeCSV returns the entire content of the CSV file present in the provided path.
func ReadWholeCSV(filePath string) ([][]string, error) {
	// Open the file for the CSV reader.
	fileReader, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s, because: %w", filePath, err)
	}
	// File reader will be closed upon return.
	defer func() { _ = fileReader.Close() }()

	// Instantiate a new CSV reader to read the given file.
	csvReader := csv.NewReader(fileReader)
	// This var will hold all the CSV content.
	var csvContent [][]string
	// Infinite loop to read the whole CSV. It terminates when the file ends.
	for {
		// Read line by line.
		content, err := csvReader.Read()
		// Ignore some errors because bank statements are ill-formatted and hence trigger a lot of false alarms.
		if err != nil && !errors.Is(err, csv.ErrFieldCount) && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("failed to read csv content: %w", err)
		}

		// This means that the file has ended.
		if len(content) == 0 {
			break
		}

		// Collect results.
		csvContent = append(csvContent, content)
	}

	return csvContent, nil
}

// WriteJSONFile writes the provided data into the JSON file at the given path.
func WriteJSONFile(data interface{}, filePath string) error {
	// Marshal the transaction list to write into file.
	jsonBytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	// Write the output file.
	if err := os.WriteFile(filePath, jsonBytes, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}
