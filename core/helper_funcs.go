package core

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/shivanshkc/ledgerconv/core/models"
)

// showDirs returns all directories present within the given directory.
// It ignores files.
func showDirs(directory string) ([]string, error) {
	// Read the given directory to get its children.
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %s, because: %w", directory, err)
	}

	// This slice will hold the directories to be returned.
	var dirs []string
	// Filter out files from the entry list.
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}

// showFiles returns all files present within the given directory.
// It ignores directories.
func showFiles(directory string) ([]string, error) {
	// Read the given directory to get its children.
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %s, because: %w", directory, err)
	}

	// This slice will hold the files to be returned.
	var files []string
	// Filter out directories from the entry list.
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// readCSV reads the provided CSV file.
func readCSV(pathToFile string) ([][]string, error) {
	// Open the file for the CSV reader.
	fileReader, err := os.Open(pathToFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s, because: %w", pathToFile, err)
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
		// Ignore some errors because bank statements are ill-formatted and hence trigger a lot of alarms.
		if err != nil && !errors.Is(err, csv.ErrFieldCount) && !errors.Is(err, io.EOF) {
			panic("Failed to read csv content: " + err.Error())
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

// genConvertedTxChecksum provides the checksum of the given transaction.
func genConvertedTxChecksum(tx *models.ConvertedTransactionDoc) (string, error) {
	// Convert to byte slice to write to hash.
	txBytes, err := json.Marshal(tx)
	if err != nil {
		return "", fmt.Errorf("failed to marshal transaction: %+v, because: %w", tx, err)
	}
	// Calculate, format and return the checksum.
	return fmt.Sprintf("%x", sha256.Sum256(txBytes)), nil
}

// prettyPrintJSON prints a prettified json.
func prettyPrintJSON(jsonLike interface{}) error {
	jsonBytes, err := json.MarshalIndent(jsonLike, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	fmt.Println(string(jsonBytes))
	return nil
}
