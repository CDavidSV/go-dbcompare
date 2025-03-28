package generate

import (
	"fmt"

	"github.com/CDavidSV/db-comparation-tool/internal/config"
	"github.com/spf13/cobra"
)

var GenerateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Generate configuration file or DSN",
	Long:    "The generate command allows generating configuration files or DSN strings for connecting to the database.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.ErrorStyle.Render("Error: 'generate' command requires a subcommand (dsn or config)"))
		cmd.Help()
	},
}

func init() {
	GenerateCmd.AddCommand(dsnCmd)
	GenerateCmd.AddCommand(configCmd)
}
