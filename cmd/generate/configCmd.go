package generate

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generates a config file",
	Long:  "Generates a json config file for connecting to the databases that want to be compared",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {

}
