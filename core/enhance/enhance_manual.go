package enhance

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shivanshkc/ledgerconv/core/models"
	"github.com/shivanshkc/ledgerconv/core/utils/io"

	"github.com/fatih/color"
)

// Manual enhances the given converted transaction manually by asking the user for inputs.
//
//nolint:funlen,cyclop // TODO?
func Manual(txn *models.ConvertedTransactionDoc) (*models.EnhancedTransactionDoc, error) {
	color.Cyan("-----------------------------------------------------------------")
	io.PrettyPrintConvTx(txn)
	color.Cyan("-----------------------------------------------------------------")

	// Quantity to be returned.
	enhanced := &models.EnhancedTransactionDoc{
		ConvertedTransactionDoc: txn,
		Categories:              &models.AmountPerCategory{},
	}

	// Map category names to their struct field pointers. This will be helpful while prompting the user.
	creditCategoriesMap := map[string]*float64{
		// Credit categories.
		"Salary":    &enhanced.Categories.Salary,
		"Returns":   &enhanced.Categories.Returns,
		"Misc":      &enhanced.Categories.Misc,
		"Ignorable": &enhanced.Categories.Ignorable,
	}

	// Do the same for debit categories as well.
	debitCategoriesMap := map[string]*float64{
		// Debit categories.
		"Essentials":  &enhanced.Categories.Essentials,
		"Investments": &enhanced.Categories.Investments,
		"Savings":     &enhanced.Categories.Savings,
		"Luxury":      &enhanced.Categories.Luxury,
		"Ignorable":   &enhanced.Categories.Ignorable,
	}

	// Decide on the categories for the prompt.
	promptCategories := creditCategoriesMap
	if txn.Amount < 0 {
		promptCategories = debitCategoriesMap
	}

	color.Blue("Provide amount distribution among categories...")
	for {
		// This will hold the total amount sum for all categories. This should be equal to the transaction amount.
		var catAmountSum float64

		// Loop over categories to take user input.
		for catName, catPtr := range promptCategories {
			// Prompt for the category amount.
			value, err := io.Prompt(fmt.Sprintf("%s component?: ", catName))
			if err != nil {
				return nil, fmt.Errorf("failed to read user input: %w", err)
			}

			// Allow users to provide empty values.
			if value == "" {
				value = "0"
			}

			// Parse the string to float.
			valueFloat, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse float input: %w", err)
			}

			// Persist in the AmountPerCategory struct.
			*catPtr = valueFloat
			catAmountSum += valueFloat
		}

		// If the given amounts sum to the main amount, we can move further.
		if catAmountSum == txn.Amount {
			break
		}
		// Show error message if amount sum does not equal the transaction amount.
		fmt.Printf("Sum of these amounts should equal the transaction amount. But %f != %f\n",
			catAmountSum, txn.Amount)
	}

	color.Cyan("-----------------------------------------------------------------")
	color.Blue("Other information...")
	// Read labels.
	commaSepLabels, err := io.Prompt("Any labels? (comma-separated, case-insensitive): ")
	if err != nil {
		return nil, fmt.Errorf("failed to read user input: %w", err)
	}

	// Parse and format labels.
	labels := strings.Split(commaSepLabels, ",")
	for i := range labels {
		labels[i] = strings.ToLower(strings.TrimSpace(labels[i]))
	}

	color.Cyan("-----------------------------------------------------------------")
	color.Yellow("Summary?: ")

	// Read summary.
	summary, err := io.Prompt("Summary: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read user input: %w", err)
	}

	// Attach information.
	enhanced.Labels = labels
	enhanced.Summary = summary

	return enhanced, nil
}
