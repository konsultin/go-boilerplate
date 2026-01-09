package coreSql

import (
	"github.com/jmoiron/sqlx"
	"github.com/konsultin/sqlk"
	"github.com/konsultin/sqlk/pq/query"
)

type ClientAuth struct {
	FindByClientId *sqlx.Stmt
}

func NewClientAuth(db *sqlk.DatabaseContext) *ClientAuth {
	return &ClientAuth{
		FindByClientId: db.MustPrepareRebind(query.Select(query.Column("*")).
			From(ClientAuthSchema).
			Where(query.Equal(query.Column("clientId"))).
			Build()),
	}
}
