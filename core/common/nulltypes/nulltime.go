package nulltypes

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NullTime es un wrapper para manejar campos de tiempo que pueden ser NULL
type NullTime struct {
	Valid bool
	Time  time.Time
}

func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		nt.Time = time.Time{}
		return nil
	}

	nt.Valid = true
	switch v := value.(type) {
	case time.Time:
		nt.Time = v
	case string:
		// Intentar parsear el string como fecha
		parsedTime, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			// Intentar con otro formato común
			parsedTime, err = time.Parse(time.RFC3339, v)
			if err != nil {
				nt.Valid = false
				return err
			}
		}
		nt.Time = parsedTime
	case []byte:
		// Convertir bytes a string y parsear
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			parsedTime, err = time.Parse(time.RFC3339, string(v))
			if err != nil {
				nt.Valid = false
				return err
			}
		}
		nt.Time = parsedTime
	default:
		nt.Valid = false
		nt.Time = time.Time{}
	}
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Time)
	}
	return json.Marshal(nil)
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		nt.Time = time.Time{}
		return nil
	}
	nt.Valid = true
	return json.Unmarshal(data, &nt.Time)
}

// ToTimePtr convierte NullTime a *time.Time
func (nt NullTime) ToTimePtr() *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
