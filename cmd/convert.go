package cmd

import (
	"fmt"
	"os"

	"github.com/shivanshkc/ledgerconv/core"

	"github.com/spf13/cobra"
)

var convertParamOutputDir string

// convertCmd represents the convert command.
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert CSV bank statements into JSON",
	Run: func(cmd *cobra.Command, args []string) {
		inputDir := args[0]
		// Core call.
		if err := core.Convert(cmd.Context(), inputDir, convertParamOutputDir); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to convert statements: %+v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&convertParamOutputDir, "output", "o", "./output",
		"The directory where the converted statement will be stored.")
}
