package sqlfmt

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NilString sql.NullString
type NilTime sql.NullTime

var jsonNull = []byte("null")

func (n NilString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return jsonNull, nil
}

func (n NilTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return jsonNull, nil
}

func NewNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  len(s) > 0,
	}
}

func NewNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

func NewNilString(s string) NilString {
	return NilString{
		String: s,
		Valid:  len(s) > 0,
	}
}

func NewNilTime(t time.Time) NilTime {
	return NilTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}
