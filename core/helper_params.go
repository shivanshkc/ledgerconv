package core

const (
	// enhancedFilename is the name of the file in which the enhanced transactions will be written.
	enhancedFilename = "enhanced-transactions.json"

	// convertedFilename is the name of the file in which the converted transactions will be written.
	convertedFilename = "converted-transactions.json"
)

const (
	// Credit categories constants.
	categorySalary  = "salary"
	categoryReturns = "returns"
	categoryMisc    = "misc"

	// Debit categories constants.
	categoryEssentials  = "essentials"
	categoryInvestments = "investments"
	categorySavings     = "savings"
	categoryLuxury      = "luxury"

	// Common categories constants.
	categoryIgnorable = "ignorable"
)

var (
	// Slices for conveniently looping over categories.
	creditCats = []string{categorySalary, categoryReturns, categoryMisc, categoryIgnorable}
	debitCats  = []string{categoryEssentials, categoryInvestments, categorySavings, categoryLuxury, categoryIgnorable}
)
