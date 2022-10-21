package cmd

import (
	"os"

	"github.com/shivanshkc/ledgerconv/core"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	enhParamOutput string
	enhParamSpec   string
)

// enhanceCmd represents the enhance command.
var enhanceCmd = &cobra.Command{
	Use:   "enhance <input-file>",
	Short: "Enhance formatted statements with custom fields.",

	// At least one argument is required.
	Args: cobra.MinimumNArgs(1),

	// Command action.
	Run: func(cmd *cobra.Command, args []string) {
		ctx, inputFile := cmd.Context(), args[0]
		// Core call.
		if err := core.NewEnhancer().Enhance(ctx, inputFile, enhParamOutput, enhParamSpec); err != nil {
			_, _ = color.New(color.FgRed).Fprintf(os.Stderr, "Failed to enhance statements: %+v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(enhanceCmd)

	enhanceCmd.Flags().StringVarP(&enhParamOutput, "output", "o", "./enhanced.json",
		"Path where the enhanced statement file will be created or updated.")

	enhanceCmd.Flags().StringVarP(&enhParamSpec, "auto-enhance-spec", "s", "",
		"Path to the auto-enhance specification file.")
}
