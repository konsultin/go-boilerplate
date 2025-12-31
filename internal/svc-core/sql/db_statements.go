package coreSql

import "github.com/konsultin/project-goes-here/libs/sqlk"

type Statements struct {
	User *User
}

func New(db *sqlk.DatabaseContext) *Statements {
	return &Statements{
		User: NewUser(db),
	}
}
