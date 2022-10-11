package cmd

import (
	"fmt"
	"os"

	"github.com/shivanshkc/ledgerconv/core"

	"github.com/spf13/cobra"
)

var enhanceParamOutputDir string

// enhanceCmd represents the enhance command.
var enhanceCmd = &cobra.Command{
	Use:   "enhance <input-dir>",
	Short: "Enhance formatted statements with custom fields.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputDir := args[0]
		// Core call.
		if err := core.Enhance(cmd.Context(), inputDir, enhanceParamOutputDir); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to enhance statements: %+v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(enhanceCmd)

	enhanceCmd.Flags().StringVarP(&enhanceParamOutputDir, "output", "o", "./output",
		"The directory where the enhanced statement will be stored.")
}
