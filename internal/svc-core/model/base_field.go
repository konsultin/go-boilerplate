package model

import (
	"encoding/json"
	"time"

	"github.com/Konsultin/project-goes-here/dto"
	"github.com/Konsultin/project-goes-here/libs/sqlk"
)

type BaseField struct {
	CreatedAt  time.Time       `db:"createdAt"`
	UpdatedAt  time.Time       `db:"updatedAt"`
	ModifiedBy *Subject        `db:"modifiedBy"`
	Version    int64           `db:"version"`
	Metadata   json.RawMessage `db:"metadata"`
}

func NewBaseField(subject *dto.Subject) BaseField {
	t := time.Now()
	return BaseField{
		CreatedAt:  t,
		UpdatedAt:  t,
		ModifiedBy: NewSubject(subject),
		Version:    1,
		Metadata:   sqlk.EmptyObjectJSON,
	}
}
