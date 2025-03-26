package internal

import (
	"encoding/json"
	"fmt"
	"math"
)

func GenerateCSVData(differences ComparisonResult, DB1Name, DB2Name string) ([][]string, error) {
	data := [][]string{
		{fmt.Sprintf("Database 1 (%s)", DB1Name), fmt.Sprintf("Database 2 (%s)", DB2Name), "", "", fmt.Sprintf("Tables missing in %s", DB1Name), fmt.Sprintf("Tables missing in %s", DB2Name)},
	}

	for i, v := range differences.DifferencesResult.DB1 {
		v2 := differences.DifferencesResult.DB2[i]

		jsonStringDB1, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			return nil, err
		}

		jsonStringDB2, err := json.MarshalIndent(v2, "", "    ")
		if err != nil {
			return nil, err
		}

		comparisonData := []string{string(jsonStringDB1), string(jsonStringDB2)}

		data = append(data, comparisonData)
	}

	// Get the length of the longest slice
	length := int(math.Max(float64(len(differences.MissingTablesInDB1)), float64(len(differences.MissingTablesInDB2))))

	dataIndex := 1
	for i := range length {
		tableDB1 := differences.MissingTablesInDB1[i]
		tableDB2 := ""

		if i <= len(differences.MissingTablesInDB2) {
			tableDB2 = differences.MissingTablesInDB2[i]
		}

		// First determine if the index is null
		if dataIndex >= len(data) {
			data = append(data, []string{"", "", "", "", tableDB1, tableDB2})
		} else {
			data[dataIndex] = append(data[dataIndex], "", "", tableDB1, tableDB2)
		}

		dataIndex++
	}

	return data, nil
}
