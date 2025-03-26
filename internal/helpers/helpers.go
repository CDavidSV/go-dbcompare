package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/CDavidSV/db-comparation-tool/internal"
	_ "github.com/lib/pq"
)

var (
	clearLine          = "\033[u\033[K"
	saveCursorPosition = "\033[s"
)

func ClearLine() {
	fmt.Print(clearLine)
}

func SaveCursorPosition() {
	fmt.Print(saveCursorPosition)
}

func getDataSourceName(driverName string, dbConfig internal.DBConfig) string {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s", driverName, dbConfig.Username, dbConfig.Password, dbConfig.HostName, dbConfig.Port, dbConfig.Database)

	if dbConfig.Params == nil {
		return dsn
	}

	if len(dbConfig.Params) == 0 {
		return dsn
	}

	params := url.Values{}
	for param, val := range dbConfig.Params {
		params.Add(param, val)
	}
	dsn += "?" + params.Encode()

	return dsn
}

func ConnectDB(driverName string, dbConfig internal.DBConfig) (*sql.DB, error) {
	dsn := getDataSourceName(driverName, dbConfig)
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
