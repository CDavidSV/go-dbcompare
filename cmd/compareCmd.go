package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/CDavidSV/db-comparation-tool/internal"
	"github.com/CDavidSV/db-comparation-tool/internal/config"
	"github.com/CDavidSV/db-comparation-tool/internal/helpers"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

var compareCmd = cobra.Command{
	Use:   "compare",
	Short: "Runs comparison between two databases",
	Long:  "Connects to the databases and compares tables found in each one.",
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath, _ := cmd.Flags().GetString("config")
		outputPath, _ := cmd.Flags().GetString("output")
		name, _ := cmd.Flags().GetString("name")

		conf, err := helpers.LoadConfigurationFile(configFilePath)
		if err != nil {
			fmt.Println(config.ErrorStyle.Render("Error reading configuration file:"), err)
			os.Exit(1)
		}

		helpers.SaveCursorPosition()

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = fmt.Sprintf(config.InfoStyle.Render(" Connecting to %s"), conf.DB1.Name)
		s.Start()

		// Connect to databases
		DB1, err := helpers.ConnectDB("postgres", conf.DB1)
		s.Stop()
		if err != nil {
			helpers.ClearLine()
			fmt.Printf(config.ErrorStyle.Render("Error connecting to %s\nError: %s\n"), conf.DB1.Name, err)
			os.Exit(1)
		}
		defer DB1.Close()

		helpers.ClearLine()
		fmt.Printf(config.SuccessStyle.Render("Connected to %s\n"), conf.DB1.Name)
		helpers.SaveCursorPosition()

		s.Suffix = fmt.Sprintf(config.InfoStyle.Render(" Connecting to %s"), conf.DB2.Name)
		s.Start()

		DB2, err := helpers.ConnectDB("postgres", conf.DB2)
		s.Stop()
		if err != nil {
			helpers.ClearLine()
			fmt.Printf(config.ErrorStyle.Render("Error connecting to %s\nError: %s\n"), conf.DB2.Name, err)
			os.Exit(1)
		}
		defer DB2.Close()

		helpers.ClearLine()
		fmt.Printf(config.SuccessStyle.Render("Connected to %s\n\n"), conf.DB2.Name)
		helpers.SaveCursorPosition()

		s.Suffix = config.InfoStyle.Render(" Running comparison")
		s.Start()

		f, err := runDBComparison(conf.DB1.Name, DB1, conf.DB2.Name, DB2)
		s.Stop()
		if err != nil {
			helpers.ClearLine()
			fmt.Printf(config.ErrorStyle.Render("Error running database comparison: %s\n"), err)
			os.Exit(1)
		}
		defer f.Close()

		helpers.ClearLine()
		fmt.Println(config.SuccessStyle.Render("✔ Comparison finished"))

		if name == "" {
			outputPath += "Comparison_Result.xlsx"
		} else {
			outputPath += name + ".xlsx"
		}

		file, err := os.Create(outputPath)
		if err != nil {
			fmt.Println(config.ErrorStyle.Render("Error saving result file:"), err)
			os.Exit(1)
		}

		_, err = f.WriteTo(file)
		if err != nil {
			fmt.Println(config.ErrorStyle.Render("Error saving result file:"), err)
			os.Exit(1)
		}

		file.Close()
		fmt.Println(config.SuccessStyle.Render("✔ Result file saved successfully"))
	},
}

func runDBComparison(DB1Name string, DB1 *sql.DB, DB2Name string, DB2 *sql.DB) (*excelize.File, error) {
	result, err := internal.CompareDatabase(DB1, DB2)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheetName := "Database Comparison Result"

	f.SetSheetName(f.GetSheetName(0), sheetName)

	style1, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})

	errorStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#f4cccc"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})

	borderStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})

	// Set titles
	f.SetCellValue(sheetName, "A1", "Num")
	f.SetCellValue(sheetName, "B1", fmt.Sprintf("Database 1 (%s)", DB1Name))
	f.SetCellValue(sheetName, "C1", fmt.Sprintf("Database 2 (%s)", DB2Name))
	f.SetCellStyle(sheetName, "A1", "C1", borderStyle)
	f.SetCellValue(sheetName, "F1", fmt.Sprintf("Tables missing in %s", DB1Name))
	f.SetCellValue(sheetName, "G1", fmt.Sprintf("Tables missing in %s", DB2Name))
	f.SetCellStyle(sheetName, "F1", "G1", borderStyle)

	// Fill in the table data
	firstCell := 2
	for i, DB1Val := range result.DifferencesResult.DB1 {

		cellNameStart, _ := excelize.CoordinatesToCellName(1, firstCell)
		cellNameEnd, _ := excelize.CoordinatesToCellName(1, firstCell+6)
		f.SetCellStyle(sheetName, cellNameStart, cellNameEnd, style1)
		f.MergeCell(sheetName, cellNameStart, cellNameEnd)

		f.SetCellValue(sheetName, cellNameStart, i+1)

		DB2Val := result.DifferencesResult.DB2[i]

		DB1ValSlice := []string{
			fmt.Sprintf("Table Name: %s", DB1Val.TableName),
			fmt.Sprintf("Column Name: %s", DB1Val.ColumnName),
			fmt.Sprintf("Data Type: %s", DB1Val.DataType),
			fmt.Sprintf("Column Default: %s", DB1Val.ColumnDefault),
			fmt.Sprintf("Is Nullable: %s", DB1Val.IsNullable),
			fmt.Sprintf("Char Max Len: %d", DB1Val.CharMaxLen),
			fmt.Sprintf("Numeric Precision: %d", DB1Val.NumericPrecision),
		}
		DB2ValSlice := []string{
			fmt.Sprintf("Table Name: %s", DB2Val.TableName),
			fmt.Sprintf("Column Name: %s", DB2Val.ColumnName),
			fmt.Sprintf("Data Type: %s", DB2Val.DataType),
			fmt.Sprintf("Column Default: %s", DB2Val.ColumnDefault),
			fmt.Sprintf("Is Nullable: %s", DB2Val.IsNullable),
			fmt.Sprintf("Char Max Len: %d", DB2Val.CharMaxLen),
			fmt.Sprintf("Numeric Precision: %d", DB2Val.NumericPrecision),
		}

		// Compare and add to cell
		for j, v1 := range DB1ValSlice {
			v2 := DB2ValSlice[j]

			cellNameDB1, _ := excelize.CoordinatesToCellName(2, firstCell+j)
			cellNameDB2, _ := excelize.CoordinatesToCellName(3, firstCell+j)

			f.SetCellValue(sheetName, cellNameDB1, v1)
			f.SetCellValue(sheetName, cellNameDB2, v2)

			if v1 != v2 {
				f.SetCellStyle(sheetName, cellNameDB1, cellNameDB2, errorStyle)
			} else {
				f.SetCellStyle(sheetName, cellNameDB1, cellNameDB2, borderStyle)
			}

		}

		firstCell += 7
	}

	// Add missing tables
	for i, v := range result.MissingTablesInDB1 {
		cellNameDB1, _ := excelize.CoordinatesToCellName(6, i+2)
		f.SetCellValue(sheetName, cellNameDB1, v)
		f.SetCellStyle(sheetName, cellNameDB1, cellNameDB1, borderStyle)
	}

	for i, v := range result.MissingTablesInDB2 {
		cellNameDB2, _ := excelize.CoordinatesToCellName(7, i+2)
		f.SetCellValue(sheetName, cellNameDB2, v)
		f.SetCellStyle(sheetName, cellNameDB2, cellNameDB2, borderStyle)
	}

	f.SetColWidth(sheetName, "B", "C", 70)
	f.SetColWidth(sheetName, "F", "G", 60)

	return f, nil
}

func init() {
	compareCmd.Flags().StringP("config", "c", "./db-compare-config.json", "path for the configuration file")
	compareCmd.Flags().StringP("output", "o", "./", "path where the comparison result file is saved")
	compareCmd.Flags().StringP("name", "n", "", "name of the comparison result file")
}
