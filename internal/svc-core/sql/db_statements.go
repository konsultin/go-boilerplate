package coreSql

import "github.com/konsultin/project-goes-here/libs/sqlk"

type Statements struct {
	User        *User
	ClientAuth  *ClientAuth
	Role        *Role
	AuthSession *AuthSessionSql
}

func New(db *sqlk.DatabaseContext) *Statements {
	return &Statements{
		User:        NewUser(db),
		ClientAuth:  NewClientAuth(db),
		Role:        NewRole(db),
		AuthSession: NewAuthSession(db),
	}
}
