package coreSql

import (
	"github.com/jmoiron/sqlx"
	"github.com/konsultin/project-goes-here/libs/sqlk"
	"github.com/konsultin/project-goes-here/libs/sqlk/pq/query"
)

type User struct {
	GetUserByXid *sqlx.Stmt
}

func NewUser(db *sqlk.DatabaseContext) *User {
	return &User{
		GetUserByXid: db.MustPrepareRebind(
			query.Select(
				query.Column("*"),
			).
				From(UserSchema).
				Where(
					query.Equal(query.Column("xid")),
				).Build(),
		),
	}
}
