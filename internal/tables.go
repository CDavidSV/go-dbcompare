package internal

import (
	"database/sql"
)

type ColumnData struct {
	TableName        string
	ColumnName       string
	DataType         NullString
	ColumnDefault    NullString
	IsNullable       string
	CharMaxLen       NullInt
	NumericPrecision NullInt
}

type Differences struct {
	DB1 []ColumnData
	DB2 []ColumnData
}

type ComparisonResult struct {
	MissingTablesInDB1 []string
	MissingTablesInDB2 []string
	DifferencesResult  Differences
}

func GetDBTableData(db *sql.DB) (map[string]map[string]ColumnData, error) {
	rows, err := db.Query(`SELECT
		c.table_name, c.column_name, c.data_type, c.column_default, c.is_nullable, c.character_maximum_length, c.numeric_precision
	FROM
		INFORMATION_SCHEMA.TABLES t
	INNER JOIN INFORMATION_SCHEMA.COLUMNS c ON t.table_name = c.table_name
	WHERE
		t.table_schema='public' AND t.table_type='BASE TABLE'
	ORDER BY
		c.table_name ASC`)

	if err != nil {
		return nil, err
	}

	// table > column > columnData
	tables := map[string]map[string]ColumnData{}

	for rows.Next() {
		var col ColumnData

		err = rows.Scan(&col.TableName, &col.ColumnName, &col.DataType, &col.ColumnDefault, &col.IsNullable, &col.CharMaxLen, &col.NumericPrecision)
		if err != nil {
			return nil, err
		}

		if _, ok := tables[col.TableName]; !ok {
			tables[col.TableName] = map[string]ColumnData{
				col.ColumnName: col,
			}
			continue
		}

		tables[col.TableName][col.ColumnName] = col
	}

	return tables, nil
}

func CompareTableCols(DB1Cols, DB2Cols map[string]ColumnData, differences Differences) Differences {
	for key, DB1Value := range DB1Cols {
		DB2Value, ok := DB2Cols[key]
		if !ok {
			differences.DB1 = append(differences.DB1, DB1Value)
			differences.DB2 = append(differences.DB2, ColumnData{
				TableName:        DB1Value.TableName,
				ColumnName:       "Null",
				DataType:         "Null",
				ColumnDefault:    "Null",
				IsNullable:       "YES",
				CharMaxLen:       0,
				NumericPrecision: 0,
			})
			continue
		}

		DB1ColVals := []any{
			DB1Value.ColumnName,
			DB1Value.CharMaxLen,
			DB1Value.ColumnDefault,
			DB1Value.DataType,
			DB1Value.IsNullable,
			DB1Value.NumericPrecision,
		}
		DB2ColVals := []any{
			DB2Value.ColumnName,
			DB2Value.CharMaxLen,
			DB2Value.ColumnDefault,
			DB2Value.DataType,
			DB2Value.IsNullable,
			DB2Value.NumericPrecision,
		}

		for i, v := range DB1ColVals {
			if v != DB2ColVals[i] {
				differences.DB1 = append(differences.DB1, DB1Value)
				differences.DB2 = append(differences.DB2, DB2Value)
				break
			}
		}

	}

	// Second loop: Check keys in DB2Cols against DB1Cols
	for key, DB2Value := range DB2Cols {
		if _, ok := DB1Cols[key]; !ok {
			differences.DB1 = append(differences.DB1, ColumnData{
				TableName:        DB2Value.TableName,
				ColumnName:       "Null",
				DataType:         "Null",
				ColumnDefault:    "Null",
				IsNullable:       "YES",
				CharMaxLen:       0,
				NumericPrecision: 0,
			})
			differences.DB2 = append(differences.DB2, DB2Value)
		}
	}

	return differences
}

func CompareDatabase(DB1 *sql.DB, DB2 *sql.DB) (ComparisonResult, error) {
	differences := Differences{
		DB1: []ColumnData{},
		DB2: []ColumnData{},
	}

	comparisonResult := ComparisonResult{
		MissingTablesInDB1: []string{},
		MissingTablesInDB2: []string{},
	}

	Database1TableData, err := GetDBTableData(DB1)
	if err != nil {
		return comparisonResult, err
	}

	Database2TableData, err := GetDBTableData(DB2)
	if err != nil {
		return comparisonResult, err
	}

	for DB1Key, DB1Value := range Database1TableData {
		DB2Value, ok := Database2TableData[DB1Key]
		if !ok {
			// Table not found in database 2
			comparisonResult.MissingTablesInDB2 = append(comparisonResult.MissingTablesInDB2, DB1Key)
			continue
		}

		differences = CompareTableCols(DB1Value, DB2Value, differences)
	}

	// Find all missing tables in database 1
	for DB2Key := range Database2TableData {
		if _, ok := Database1TableData[DB2Key]; !ok {
			comparisonResult.MissingTablesInDB1 = append(comparisonResult.MissingTablesInDB1, DB2Key)
		}
	}

	comparisonResult.DifferencesResult = differences

	return comparisonResult, nil
}
