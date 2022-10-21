package cmd

import (
	"os"

	"github.com/shivanshkc/ledgerconv/core"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	enhParamOutput   string
	enhParamSpec     string
	enhParamOnlyAuto bool
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
		err := core.NewEnhancer().Enhance(ctx, inputFile, enhParamOutput, enhParamSpec, enhParamOnlyAuto)
		if err != nil {
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

	enhanceCmd.Flags().BoolVar(&enhParamOnlyAuto, "only-auto", false,
		"Only auto-enhance transactions. Skip the manual ones.")
}
