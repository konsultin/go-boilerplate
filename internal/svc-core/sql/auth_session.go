package coreSql

import (
	"github.com/jmoiron/sqlx"
	"github.com/konsultin/project-goes-here/libs/sqlk"
	"github.com/konsultin/project-goes-here/libs/sqlk/pq/query"
)

type AuthSessionSql struct {
	FindByXid  *sqlx.Stmt
	DeleteById *sqlx.Stmt
	Insert     *sqlx.NamedStmt
}

func NewAuthSession(db *sqlk.DatabaseContext) *AuthSessionSql {
	return &AuthSessionSql{
		FindByXid: db.MustPrepareRebind(
			query.Select(
				query.Column("*"),
			).
				From(AuthSessionSchema).
				Where(
					query.Equal(query.Column("xid")),
				).Build(),
		),
		DeleteById: db.MustPrepareRebind(
			query.Delete(AuthSessionSchema).
				Where(
					query.Equal(query.Column("id")),
				).Build(),
		),
		Insert: db.MustPrepareNamed(
			query.Insert(AuthSessionSchema, query.AllColumns).Build(),
		),
	}
}
