package internal

import "fmt"

type NullString string

func (ns *NullString) Scan(value any) error {
	if value == nil {
		*ns = "Null"
		return nil
	}

	switch v := value.(type) {
	case string:
		*ns = NullString(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type %T", value)
	}
}

type NullInt int

func (ni *NullInt) Scan(value any) error {
	if value == nil {
		*ni = 0
		return nil
	}

	switch v := value.(type) {
	case int64:
		*ni = NullInt(v)
		return nil
	case int:
		*ni = NullInt(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type %T", value)
	}
}

type DBConfig struct {
	Name     string            `json:"name"`
	HostName string            `json:"host"`
	Port     uint16            `json:"port"`
	Database string            `json:"database"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Params   map[string]string `json:"params,omitempty"`
}

type Configuration struct {
	DB1 DBConfig `json:"database1"`
	DB2 DBConfig `json:"database2"`
}
