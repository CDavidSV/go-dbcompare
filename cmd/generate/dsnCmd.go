package generate

import (
	"fmt"
	"os"
	"strings"

	"github.com/CDavidSV/db-comparation-tool/internal/config"
	"github.com/CDavidSV/db-comparation-tool/internal/helpers"
	"github.com/spf13/cobra"
)

var dsnCmd = &cobra.Command{
	Use:   "dsn",
	Short: "Generates a DSN string",
	Long:  "Generates a Data Source Name connection string",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		database, _ := cmd.Flags().GetString("database")
		port, _ := cmd.Flags().GetUint16("port")
		params, _ := cmd.Flags().GetStringArray("param")

		paramsMap := make(map[string]string)
		for _, param := range params {
			sep := strings.Split(param, "=")

			if len(sep) != 2 {
				fmt.Println(config.ErrorStyle.Render("Error: param value must be in the following format key=value. Got: ", param))
				os.Exit(1)
			}

			paramsMap[sep[0]] = sep[1]
		}

		dsn := helpers.GetDataSourceName("postgres", user, password, host, database, port, paramsMap)

		fmt.Printf("DSN: %q\n", dsn)
	},
}

func init() {
	dsnCmd.Flags().StringP("host", "s", "", "Database server hostname or IP address")
	dsnCmd.Flags().StringP("user", "u", "", "Username for database authentication")
	dsnCmd.Flags().StringP("password", "a", "", "Password for database authentication")
	dsnCmd.Flags().StringP("database", "d", "", "Name of the database to connect to")
	dsnCmd.Flags().Uint16P("port", "p", 5432, "Port number for the database connection")
	dsnCmd.Flags().StringArrayP("param", "e", []string{}, "Optional connection parameters (e.g., --param key=value)")

	dsnCmd.MarkFlagRequired("host")
	dsnCmd.MarkFlagRequired("user")
	dsnCmd.MarkFlagRequired("password")
	dsnCmd.MarkFlagRequired("database")
}
