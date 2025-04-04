package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/CDavidSV/go-dbcompare/internal"
	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

var (
	clearLine          = "\033[u\033[K\r"
	saveCursorPosition = "\033[s"
)

func ClearLine() {
	fmt.Print(clearLine)
}

func SaveCursorPosition() {
	fmt.Print(saveCursorPosition)
}

func GetDataSourceName(driverName string, username, password, host, database string, port uint16, connParams map[string]string) string {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s", driverName, username, password, host, port, database)

	if connParams == nil {
		return dsn
	}

	if len(connParams) == 0 {
		return dsn
	}

	params := url.Values{}
	for param, val := range connParams {
		params.Add(param, val)
	}
	dsn += "?" + params.Encode()

	return dsn
}

func ConnectDB(driverName string, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func LoadConfigurationFile(path string) (internal.Configuration, error) {
	conf := internal.Configuration{
		DB1: internal.DBConfig{
			Name: "DB1",
		},
		DB2: internal.DBConfig{
			Name: "DB2",
		},
	}

	if !strings.HasSuffix(path, ".json") {
		return conf, fmt.Errorf("configuration file must be in json format")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return conf, err
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

func SaveAsExcel(result internal.ComparisonResult, DB1Name string, DB2Name string, output string) error {
	f := excelize.NewFile()
	defer f.Close()
	sheetName := "Table Comparison"

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
	f.SetColWidth(sheetName, "B", "C", 70)

	// Create a new sheet to show missing tables in each database
	sheetName = "Missing tables"
	f.NewSheet(sheetName)

	f.SetCellValue(sheetName, "A1", fmt.Sprintf("Tables missing in %s (Present in %s)", DB1Name, DB2Name))
	f.SetCellValue(sheetName, "B1", fmt.Sprintf("Tables missing in %s (Present in %s)", DB2Name, DB1Name))
	f.SetCellStyle(sheetName, "A1", "B1", borderStyle)

	// Add missing tables
	for i, v := range result.MissingTablesInDB1 {
		cellNameDB1, _ := excelize.CoordinatesToCellName(1, i+2)
		f.SetCellValue(sheetName, cellNameDB1, v)
		f.SetCellStyle(sheetName, cellNameDB1, cellNameDB1, borderStyle)
	}

	for i, v := range result.MissingTablesInDB2 {
		cellNameDB2, _ := excelize.CoordinatesToCellName(2, i+2)
		f.SetCellValue(sheetName, cellNameDB2, v)
		f.SetCellStyle(sheetName, cellNameDB2, cellNameDB2, borderStyle)
	}

	f.SetColWidth(sheetName, "A", "B", 70)

	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = f.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}
