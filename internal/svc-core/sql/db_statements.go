package coreSql

import "github.com/konsultin/sqlk"

type Statements struct {
	User           *User
	UserCredential *UserCredentialSql
	ClientAuth     *ClientAuth
	Role           *Role
	AuthSession    *AuthSessionSql
}

func New(db *sqlk.DatabaseContext) *Statements {
	return &Statements{
		User:           NewUser(db),
		UserCredential: NewUserCredential(db),
		ClientAuth:     NewClientAuth(db),
		Role:           NewRole(db),
		AuthSession:    NewAuthSession(db),
	}
}
