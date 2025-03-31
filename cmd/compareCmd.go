package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/CDavidSV/db-comparation-tool/internal"
	"github.com/CDavidSV/db-comparation-tool/internal/config"
	"github.com/CDavidSV/db-comparation-tool/internal/helpers"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:     "compare",
	Aliases: []string{"c"},
	Short:   "Runs comparison between two databases",
	Long:    "Connects to the databases and compares tables found in each one.",
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath, _ := cmd.Flags().GetString("config")
		outputPath, _ := cmd.Flags().GetString("output")
		name, _ := cmd.Flags().GetString("name")
		dsn1, _ := cmd.Flags().GetString("dsn1")
		dsn2, _ := cmd.Flags().GetString("dsn2")

		db1Name := "DB1"
		db2Name := "DB2"
		if dsn1 == "" || dsn2 == "" {
			conf, err := helpers.LoadConfigurationFile(configFilePath)
			if err != nil {
				fmt.Println(config.ErrorStyle.Render("Error reading configuration file:"), err)
				os.Exit(1)
			}

			db1Name = conf.DB1.Name
			db2Name = conf.DB2.Name

			dsn1 = helpers.GetDataSourceName("postgres", conf.DB1.Username, conf.DB1.Password, conf.DB1.HostName, conf.DB1.Database, conf.DB1.Port, conf.DB1.Params)
			dsn2 = helpers.GetDataSourceName("postgres", conf.DB2.Username, conf.DB2.Password, conf.DB2.HostName, conf.DB2.Database, conf.DB2.Port, conf.DB2.Params)
		}

		helpers.SaveCursorPosition()

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = fmt.Sprintf(config.InfoStyle.Render(" Connecting to %s"), db1Name)
		s.Start()

		// Connect to databases
		DB1, err := helpers.ConnectDB("postgres", dsn1)
		s.Stop()
		if err != nil {
			helpers.ClearLine()
			fmt.Printf(config.ErrorStyle.Render("Error connecting to %s\nError: %s\n"), db1Name, err)
			os.Exit(1)
		}
		defer DB1.Close()

		helpers.ClearLine()
		fmt.Printf(config.SuccessStyle.Render("Connected to %s\n"), db1Name)
		helpers.SaveCursorPosition()

		s.Suffix = fmt.Sprintf(config.InfoStyle.Render(" Connecting to %s"), db2Name)
		s.Start()

		DB2, err := helpers.ConnectDB("postgres", dsn2)
		s.Stop()
		if err != nil {
			helpers.ClearLine()
			fmt.Printf(config.ErrorStyle.Render("Error connecting to %s\nError: %s\n"), db2Name, err)
			os.Exit(1)
		}
		defer DB2.Close()

		helpers.ClearLine()
		fmt.Printf(config.SuccessStyle.Render("Connected to %s\n\n"), db2Name)
		helpers.SaveCursorPosition()

		s.Suffix = config.InfoStyle.Render(" Running comparison")
		s.Start()

		result, err := internal.CompareDatabase(DB1, DB2)
		s.Stop()
		if err != nil {
			helpers.ClearLine()
			fmt.Printf(config.ErrorStyle.Render("Error running database comparison: %s\n"), err)
			os.Exit(1)
		}

		helpers.ClearLine()
		fmt.Println(config.SuccessStyle.Render("✔ Comparison finished"))

		if name == "" {
			timestamp := time.Now().Format("20060102_150405")

			outputPath += "Comparison_Result_" + timestamp + ".xlsx"
		} else {
			outputPath += name + ".xlsx"
		}

		err = helpers.SaveAsExcel(result, db1Name, db2Name, outputPath)
		if err != nil {
			fmt.Println(config.ErrorStyle.Render("Error saving result file:"), err)
			os.Exit(1)
		}

		fmt.Println(config.SuccessStyle.Render("✔ Result file saved successfully"))
	},
}

func init() {
	compareCmd.Flags().StringP("config", "c", "./db-compare-config.json", "path for the configuration file")
	compareCmd.Flags().StringP("output", "o", "./", "path where the comparison result file is saved")
	compareCmd.Flags().StringP("name", "n", "", "name of the comparison result file")
	compareCmd.Flags().String("dsn1", "", "connection string for the first database")
	compareCmd.Flags().String("dsn2", "", "connection string for the first database")
}
