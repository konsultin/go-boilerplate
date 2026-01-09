package model

import (
	"encoding/json"

	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/sqlk"
	"github.com/konsultin/timek"
)

type BaseField struct {
	CreatedAt  timek.Time      `db:"created_at"`
	UpdatedAt  timek.Time      `db:"updated_at"`
	ModifiedBy *Subject        `db:"modified_by"`
	Version    int64           `db:"version"`
	Metadata   json.RawMessage `db:"metadata"`
}

func NewBaseField(subject *dto.Subject) BaseField {
	t := timek.Now()
	return BaseField{
		CreatedAt:  t,
		UpdatedAt:  t,
		ModifiedBy: NewSubject(subject),
		Version:    1,
		Metadata:   sqlk.EmptyObjectJSON,
	}
}

// NewBaseFieldFromModel creates BaseField from model.Subject (for service layer use)
func NewBaseFieldFromModel(subject *Subject) BaseField {
	t := timek.Now()
	return BaseField{
		CreatedAt:  t,
		UpdatedAt:  t,
		ModifiedBy: subject,
		Version:    1,
		Metadata:   sqlk.EmptyObjectJSON,
	}
}
