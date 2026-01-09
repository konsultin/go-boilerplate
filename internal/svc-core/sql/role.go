package coreSql

import (
	"github.com/jmoiron/sqlx"
	"github.com/konsultin/sqlk"
	"github.com/konsultin/sqlk/pq/query"
)

type Role struct {
	FindById  *sqlx.Stmt
	FindByXid *sqlx.Stmt
}

func NewRole(db *sqlk.DatabaseContext) *Role {
	return &Role{
		FindById: db.MustPrepareRebind(query.Select(query.Column("*")).
			From(RoleSchema).
			Where(query.Equal(query.Column("id"))).
			Build()),
		FindByXid: db.MustPrepareRebind(query.Select(query.Column("*")).
			From(RoleSchema).
			Where(query.Equal(query.Column("xid"))).
			Build()),
	}
}
