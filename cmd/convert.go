package cmd

import (
	"fmt"
	"os"

	"github.com/shivanshkc/ledgerconv/core"

	"github.com/spf13/cobra"
)

var convParamOutput string

// convertCmd represents the convert command.
var convertCmd = &cobra.Command{
	Use:   "convert <input-dir>",
	Short: "Convert CSV bank statements into JSON",

	// At least one argument is required.
	Args: cobra.MinimumNArgs(1),

	// Command action.
	Run: func(cmd *cobra.Command, args []string) {
		ctx, inputDir := cmd.Context(), args[0]
		// Core call.
		if err := core.NewConverter().Convert(ctx, inputDir, convParamOutput); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to convert statements: %+v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&convParamOutput, "output", "o", "./converted.json",
		"Path where the converted statement file will be created or updated.")
}
