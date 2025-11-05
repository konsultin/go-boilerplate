package model

type User struct {
	BaseField
	Id       int64  `db:"id"`
	Xid      string `db:"xid"`
	FullName string `db:"fullName"`
}
