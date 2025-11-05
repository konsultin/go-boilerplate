package coreSql

import (
	"github.com/Konsultin/project-goes-here/libs/sqlk"
	"github.com/Konsultin/project-goes-here/libs/sqlk/pq/query"
	"github.com/jmoiron/sqlx"
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
