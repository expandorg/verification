package nulls

import (
	"database/sql"
	"encoding/json"
)

type String struct {
	sql.NullString
}

func NewString(n interface{}) String {
	if n == nil {
		return String{
			sql.NullString{
				Valid: false,
			},
		}
	}
	return String{
		sql.NullString{
			String: n.(string),
			Valid:  true,
		},
	}
}

func (n String) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

func (n *String) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		n.Valid = true
		n.String = *s
	} else {
		n.Valid = false
	}
	return nil

}
