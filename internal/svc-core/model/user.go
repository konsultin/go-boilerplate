package model

import (
	"database/sql"

	"github.com/konsultin/project-goes-here/dto"
)

type User struct {
	BaseField
	Id       int64                  `db:"id"`
	Xid      string                 `db:"xid"`
	FullName string                 `db:"full_name"`
	Phone    string                 `db:"phone"`
	Email    string                 `db:"email"`
	Age      sql.NullString         `db:"age"`
	Avatar   sql.NullString         `db:"avatar"`
	StatusId dto.ControlStatus_Enum `db:"status_id"`
}

func NewUser() *User {
	return &User{}
}
