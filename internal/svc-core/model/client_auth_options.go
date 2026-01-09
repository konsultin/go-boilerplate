package model

import (
	"encoding/json"
	"database/sql/driver"
	"github.com/konsultin/sqlk"
)

type ClientAuthOptions struct {
	ClientSecret  string `json:"clientSecret"`
	TokenLifetime int64  `json:"tokenLifetime"`
}

func (m *ClientAuthOptions) Scan(src interface{}) error {
	return sqlk.ScanJSON(src, m)
}

func (m *ClientAuthOptions) Value() (driver.Value, error) {
	return json.Marshal(m)
}
