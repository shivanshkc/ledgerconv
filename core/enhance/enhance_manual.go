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

	// Declaring categories for prompting the user.
	// TODO: This is clumsy and difficult to maintain if a new category is introduced.
	creditCategories := []string{"Salary", "Returns", "Misc", "Ignorable"}
	debitCategories := []string{"Essentials", "Investments", "Savings", "Luxury", "Ignorable"}

	// Map category names to their struct field pointers. This will be helpful while prompting the user.
	creditCategoriesMap := map[string]*float64{
		// Credit categories.
		creditCategories[0]: &enhanced.Categories.Salary,
		creditCategories[1]: &enhanced.Categories.Returns,
		creditCategories[2]: &enhanced.Categories.Misc,
		creditCategories[3]: &enhanced.Categories.Ignorable,
	}

	// Do the same for debit categories as well.
	debitCategoriesMap := map[string]*float64{
		// Debit categories.
		debitCategories[0]: &enhanced.Categories.Essentials,
		debitCategories[1]: &enhanced.Categories.Investments,
		debitCategories[2]: &enhanced.Categories.Savings,
		debitCategories[3]: &enhanced.Categories.Luxury,
		debitCategories[4]: &enhanced.Categories.Ignorable,
	}

	// Decide on the categories for the prompt.
	promptCategories, promptCategoriesMap := creditCategories, creditCategoriesMap
	if txn.Amount < 0 {
		promptCategories, promptCategoriesMap = debitCategories, debitCategoriesMap
	}

	color.Blue("Provide amount distribution among categories...")
	for {
		// This will hold the total amount sum for all categories. This should be equal to the transaction amount.
		var catAmountSum float64

		// Loop over categories to take user input.
		for _, catName := range promptCategories {
			// Get the corresponding pointer.
			catPtr := promptCategoriesMap[catName]

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
	// Read summary.
	summary, err := io.Prompt("Summary?: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read user input: %w", err)
	}

	// Attach information.
	enhanced.Labels = labels
	enhanced.Summary = summary
	enhanced.AutoEnhanced = false

	return enhanced, nil
}
