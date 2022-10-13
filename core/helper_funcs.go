package core

import (
	"bufio"
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/shivanshkc/ledgerconv/core/models"

	"github.com/fatih/color"
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

// writeJSON writes the provided JSON content into the given file.
func writeJSON(jsonLike interface{}, filepath string) error {
	// Marshal the transaction list to write into file.
	jsonBytes, err := json.MarshalIndent(jsonLike, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	// Write the output file.
	if err := os.WriteFile(filepath, jsonBytes, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// printConvertedTx prints the ConvertedTransactionDoc prettily.
func printConvertedTx(doc *models.ConvertedTransactionDoc) {
	keyColor := color.New(color.FgMagenta)
	valColor := color.New(color.FgCyan)

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	keyColor.Fprint(writer, "Account\t:\t")
	valColor.Fprint(writer, doc.AccountName)
	fmt.Fprintln(writer, "")

	keyColor.Fprint(writer, "Amount\t:\t")
	valColor.Fprint(writer, doc.Amount)
	fmt.Fprintln(writer, "")

	keyColor.Fprint(writer, "Timestamp\t:\t")
	valColor.Fprint(writer, doc.Timestamp)
	fmt.Fprintln(writer, "")

	keyColor.Fprint(writer, "Bank serial\t:\t")
	valColor.Fprint(writer, doc.BankSerial)
	fmt.Fprintln(writer, "")

	keyColor.Fprint(writer, "Payment mode\t:\t")
	valColor.Fprint(writer, doc.BankPaymentMode)
	fmt.Fprintln(writer, "")

	keyColor.Fprint(writer, "Bank remarks\t:\t")
	valColor.Fprint(writer, doc.BankRemarks)
	fmt.Fprintln(writer, "")

	writer.Flush()
}

// prompt the user for an input.
func prompt(text string) (string, error) {
	// Create a reader to read from stdin.
	reader := bufio.NewReader(os.Stdin)
	// Print the prompt text.
	_, _ = color.New(color.FgMagenta).Print(text)

	// Read user's input.
	value, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read stdin: %w", err)
	}

	// Return trimmed value.
	return strings.TrimSpace(value), nil
}
