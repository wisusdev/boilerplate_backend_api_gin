package nulltypes

import (
	"database/sql/driver"
	"encoding/json"
)

// NullString es un wrapper para sql.NullString, pero con soporte para JSON
// Puedes agregar otros tipos nulos similares si los necesitas

type NullString struct {
	Valid  bool
	String string
}

func (ns *NullString) Scan(value interface{}) error {
	ns.Valid = value != nil
	if value != nil {
		switch v := value.(type) {
		case string:
			ns.String = v
		case []byte:
			ns.String = string(v)
		case nil:
			ns.Valid = false
			ns.String = ""
		default:
			ns.String = ""
		}
	}
	return nil
}

func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		ns.String = ""
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.String)
}
