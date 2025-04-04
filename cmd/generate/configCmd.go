package generate

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/CDavidSV/go-dbcompare/internal"
	"github.com/CDavidSV/go-dbcompare/internal/config"
	"github.com/CDavidSV/go-dbcompare/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"c"},
	Short:   "Generates a config file",
	Long:    "Generates a json config file for connecting to the databases that want to be compared",
	Run: func(cmd *cobra.Command, args []string) {
		// Connection parameters for database 1
		db1Name, _ := cmd.Flags().GetString("db1-name")
		db1Host, _ := cmd.Flags().GetString("db1-host")
		db1Port, _ := cmd.Flags().GetUint16("db1-port")
		db1Database, _ := cmd.Flags().GetString("db1-database")
		db1Username, _ := cmd.Flags().GetString("db1-username")
		db1Password, _ := cmd.Flags().GetString("db1-password")
		db1Param, _ := cmd.Flags().GetStringArray("db1-param")

		// Connection parameters for database 2
		db2Name, _ := cmd.Flags().GetString("db2-name")
		db2Host, _ := cmd.Flags().GetString("db2-host")
		db2Port, _ := cmd.Flags().GetUint16("db2-port")
		db2Database, _ := cmd.Flags().GetString("db2-database")
		db2Username, _ := cmd.Flags().GetString("db2-username")
		db2Password, _ := cmd.Flags().GetString("db2-password")
		db2Param, _ := cmd.Flags().GetStringArray("db2-param")

		// Asks for the required connection parameters if not provided in flags
		askConnParamsDB(1, &db1Name, &db1Host, &db1Port, &db1Database, &db1Username, &db1Password, &db1Param)
		askConnParamsDB(2, &db2Name, &db2Host, &db2Port, &db2Database, &db2Username, &db2Password, &db2Param)

		// Generate config file
		conf := internal.Configuration{
			DB1: internal.DBConfig{
				Name:     db1Name,
				HostName: db1Host,
				Port:     db1Port,
				Database: db1Database,
				Username: db1Username,
				Password: db1Password,
			},
			DB2: internal.DBConfig{
				Name:     db2Name,
				HostName: db2Host,
				Port:     db2Port,
				Database: db2Database,
				Username: db2Username,
				Password: db2Password,
			},
		}

		configJson, _ := json.Marshal(conf)

		err := os.WriteFile("db-compare-config.json", configJson, 0644)
		if err != nil {
			fmt.Println(config.ErrorStyle.Render("Could not save file: "), err)
		}

		fmt.Println()
	},
}

func askConnParamsDB(dbID uint8, name *string, host *string, port *uint16, database *string, username *string, password *string, param *[]string) {
	var output ui.TextInputValue

	// Prompt for name
	if *name == "" {
		p := tea.NewProgram(ui.InitialTextInputModel(ui.TextInputOptions{
			Label:       fmt.Sprintf("Name to refer database %d", dbID),
			Placeholder: "Production database",
			CharLimit:   50,
			Required:    true,
		}, &output))
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}

		*name = output.Value
	}

	// Prompt for host
	if *host == "" {
		p := tea.NewProgram(ui.InitialTextInputModel(ui.TextInputOptions{
			Label:       "Database Host",
			Placeholder: "127.0.0.1",
			Required:    true,
		}, &output))
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		*host = output.Value
	}

	// Prompt for port
	if *port == 0 {
		p := tea.NewProgram(ui.InitialTextInputModel(ui.TextInputOptions{
			Label:       "Database Port",
			Placeholder: "5432",
			Required:    true,
		}, &output))
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		portValue, err := strconv.ParseUint(output.Value, 10, 16)
		if err != nil {
			log.Fatal("Invalid port number")
		}
		*port = uint16(portValue)
	}

	// Prompt for database name
	if *database == "" {
		p := tea.NewProgram(ui.InitialTextInputModel(ui.TextInputOptions{
			Label:       "Database Name",
			Placeholder: "my_database",
			Required:    true,
		}, &output))
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		*database = output.Value
	}

	// Prompt for username
	if *username == "" {
		p := tea.NewProgram(ui.InitialTextInputModel(ui.TextInputOptions{
			Label:       "Database Username",
			Placeholder: "admin",
			Required:    true,
		}, &output))
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		*username = output.Value
	}

	// Prompt for password
	if *password == "" {
		p := tea.NewProgram(ui.InitialTextInputModel(ui.TextInputOptions{
			Label:     "Database Password",
			Required:  true,
			MaskInput: true,
		}, &output))
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		*password = output.Value
	}

	// Prompt for additional parameters
	if param == nil || len(*param) == 0 {
		p := tea.NewProgram(ui.InitialTextInputModel(ui.TextInputOptions{
			Label:       "Additional Connection Parameters (comma-separated)",
			Placeholder: "sslmode=disable,connect_timeout=10",
			Required:    false,
			ValidationFunction: func(value string) error {
				// Separate the input string by coma
				separatedParams := strings.Split(value, ",")

				// Validate each connection parameter
				for _, p := range separatedParams {
					if len(strings.Split(p, "=")) != 2 {
						return fmt.Errorf("param value must be in the following format key=value. Got: %s", p)
					}
				}
				return nil
			},
		}, &output))
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}

		*param = strings.Split(output.Value, ",")
	}
}

func init() {
	configCmd.Flags().String("db1-name", "", "Name of first database")
	configCmd.Flags().String("db1-host", "", "Hostname of first database")
	configCmd.Flags().Uint16("db1-port", 0, "Port of first database")
	configCmd.Flags().String("db1-database", "", "Database name of first database")
	configCmd.Flags().String("db1-username", "", "Username for first database")
	configCmd.Flags().String("db1-password", "", "Password for first database")
	configCmd.Flags().StringArray("db1-param", []string{}, "Optional connection parameters for first database (e.g., --param key=value)")

	configCmd.Flags().String("db2-name", "", "Name of second database")
	configCmd.Flags().String("db2-host", "", "Hostname of second database")
	configCmd.Flags().Uint16("db2-port", 0, "Port of second database")
	configCmd.Flags().String("db2-database", "", "Database name of second database")
	configCmd.Flags().String("db2-username", "", "Username for second database")
	configCmd.Flags().String("db2-password", "", "Password for second database")
	configCmd.Flags().StringArray("db2-param", []string{}, "Optional connection parameters for second database (e.g., --param key=value)")
}
