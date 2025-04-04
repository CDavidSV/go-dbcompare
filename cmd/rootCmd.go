package cmd

import (
	"log"

	"github.com/CDavidSV/go-dbcompare/cmd/generate"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "db-diff",
	Short: "Program that compares two databases",
	Long:  "CLI tool that helps identify differences in tables, and exporting results to an Excel file. Useful for database migrations, audits, and integrity checks.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(compareCmd)
	rootCmd.AddCommand(generate.GenerateCmd)
}
