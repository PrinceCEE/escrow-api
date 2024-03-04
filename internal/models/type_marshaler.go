package models

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

type NullTime struct {
	sql.NullTime
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(n.String)
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(n.Time)
}
