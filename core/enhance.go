package core

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/shivanshkc/ledgerconv/core/models"
)

// enhancedFilename is the name of the file in which the enhanced transactions will be written.
const enhancedFilename = "enhanced-transactions.json"

var (
	creditCats = []string{"Salary", "Returns", "Misc", "Ignorable"}
	debitCats  = []string{"Essentials", "Investments", "Savings", "Luxury", "Ignorable"}
)

// Enhance adds the custom fields with zero values to the statement present in the input directory, and places this
// new statement in the output directory.
//
// If there is already a statement in the output directory, then conflicting statement entries are skipped.
//
// This is an idempotent operation.
//
//nolint:funlen,cyclop // Core functions are allowed to be big.
func Enhance(ctx context.Context, inputDir string, outputDir string) error {
	// Path to the existing enhanced statement file.
	enhancedFilepath := path.Join(outputDir, enhancedFilename)
	// Open the enhanced statement file to load existing enhanced transactions.
	enhancedFileReader, err := os.Open(enhancedFilepath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to open file: %s, because: %w", enhancedFilepath, err)
	}

	// Decode the enhanced statement into a slice. If the file did not exist, the reader will be nil, and no decoding
	// will take place.
	var enhancedStatement []*models.EnhancedTransactionDoc //nolint:prealloc // Cannot pre-allocate.
	if enhancedFileReader != nil {
		defer func() { _ = enhancedFileReader.Close() }()

		if err := json.NewDecoder(enhancedFileReader).Decode(&enhancedStatement); err != nil {
			return fmt.Errorf("failed to decode the enhanced statement: %w", err)
		}
	}

	// This maps the enhanced transactions to their correlation IDs.
	enhancedStatementMap := map[string]*models.EnhancedTransactionDoc{}
	for _, tx := range enhancedStatement {
		enhancedStatementMap[tx.DocCorrelationID] = tx
	}

	// Path to the converted statement file.
	convertedFilepath := path.Join(inputDir, convertedFilename)
	// Open the converted statement file to enhance the transactions.
	convertedFileReader, err := os.Open(convertedFilepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %s, because: %w", convertedFilepath, err)
	}
	defer func() { _ = convertedFileReader.Close() }()

	// Decode the converted statement into a slice.
	var convertedStatement []*models.ConvertedTransactionDoc
	if err := json.NewDecoder(convertedFileReader).Decode(&convertedStatement); err != nil {
		return fmt.Errorf("failed to decode the converted statement: %w", err)
	}

	// Loop over all converted transactions to enhance them.
	for _, txn := range convertedStatement {
		// Generate checksum.
		correlationID, err := genConvertedTxChecksum(txn)
		if err != nil {
			return fmt.Errorf("failed to generate checksum: %w", err)
		}
		// See if enhanced transaction exists.
		if _, exists := enhancedStatementMap[correlationID]; exists {
			continue
		}

		fmt.Println(">> NEW TRANSACTION:", correlationID)
		_ = prettyPrintJSON(txn)

		// Prompt the user for inputs.
		enhanced, err := takeUserInput(txn)
		if err != nil {
			return fmt.Errorf("failed to take user input: %w", err)
		}

		enhanced.ConvertedTransactionDoc = txn
		enhanced.DocCorrelationID = correlationID

		enhancedStatementMap[correlationID] = enhanced
		enhancedStatement = append(enhancedStatement, enhanced)
		break
	}

	// Sort the enhanced statements.
	sort.SliceStable(enhancedStatement, func(i, j int) bool {
		return enhancedStatement[i].Timestamp.After(enhancedStatement[j].Timestamp)
	})

	// Marshal the transaction list to write into file.
	statementBytes, err := json.MarshalIndent(enhancedStatement, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal transaction list: %w", err)
	}

	// Name of the output file.
	outputFilePath := path.Join(outputDir, enhancedFilename)
	// Write the output file.
	if err := os.WriteFile(outputFilePath, statementBytes, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// takeUserInput prompts the user for inputs required to create an enhanced transaction.
//
//nolint:funlen // TODO
func takeUserInput(txn *models.ConvertedTransactionDoc) (*models.EnhancedTransactionDoc, error) {
	// This will hold the amount-per-category info.
	amountPerCat := new(models.AmountPerCategory)

	// Mapping category names to their struct field pointers. This will be helpful while prompting the user.
	catMap := map[string]*float64{
		creditCats[0]: &amountPerCat.Salary,
		creditCats[1]: &amountPerCat.Returns,
		creditCats[2]: &amountPerCat.Misc,
		creditCats[3]: &amountPerCat.Ignorable,
		debitCats[0]:  &amountPerCat.Salary,
		debitCats[1]:  &amountPerCat.Returns,
		debitCats[2]:  &amountPerCat.Misc,
		debitCats[3]:  &amountPerCat.Luxury,
		debitCats[4]:  &amountPerCat.Ignorable,
	}

	// Decide on the categories for the prompt.
	promptCats := debitCats
	if txn.Amount > 0 {
		promptCats = creditCats
	}

	fmt.Println(">> AMOUNT PER CATEGORY INFORMATION")
	for {
		// This will hold the total amount sum for all categories. This should be equal to the transaction amount.
		var catAmountSum float64

		// Loop over categories to take user input.
		for _, cat := range promptCats {
			// Get the pointer to the struct field.
			ptr := catMap[cat]

			// Prompt for the category amount.
			value, err := prompt(cat + ": ")
			if err != nil {
				return nil, fmt.Errorf("failed to read user input: %w", err)
			}

			// Parse the string to float.
			valueFloat, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse float input: %w", err)
			}

			// Persist in the AmountPerCategory struct.
			*ptr = valueFloat
			catAmountSum += valueFloat
		}

		// If the given amounts sum to the main amount, we can move further.
		if catAmountSum == txn.Amount {
			break
		}
		fmt.Printf("Sum of these amounts should equal the transaction amount. %f != %f\n", catAmountSum, txn.Amount)
	}

	fmt.Println(">> OTHER INFORMATION")
	// Read tags.
	commaSepTags, err := prompt("Tags: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read user input: %w", err)
	}

	// Parse and format tags.
	tags := strings.Split(commaSepTags, ",")
	for i := range tags {
		tags[i] = strings.ToLower(strings.TrimSpace(tags[i]))
	}

	// Read remarks.
	remarks, err := prompt("Remarks: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read user input: %w", err)
	}

	return &models.EnhancedTransactionDoc{
		AmountPerCategory: amountPerCat,
		Tags:              tags,
		Remarks:           remarks,
	}, nil
}

// prompt the user for an input.
func prompt(text string) (string, error) {
	// Create a reader to read from stdin.
	reader := bufio.NewReader(os.Stdin)
	// Print the prompt text.
	_, _ = fmt.Fprint(os.Stdout, text)

	// Read user's input.
	value, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read stdin: %w", err)
	}

	// Return trimmed value.
	return strings.TrimSpace(value), nil
}
