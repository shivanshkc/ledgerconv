package cmd

import (
	"fmt"
	"os"

	"github.com/shivanshkc/ledgerconv/core"

	"github.com/spf13/cobra"
)

var rootParamOutputDir string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "ledgerconv <input-directory>",
	Short: "Convert CSV bank statements into JSON",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputDir := args[0]
		// Core call.
		if err := core.Convert(cmd.Context(), inputDir, rootParamOutputDir); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to convert statements: %+v\n", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&rootParamOutputDir, "output", "o", "./converted",
		"The directory where the converted statement will be stored.")
}
