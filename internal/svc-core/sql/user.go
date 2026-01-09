package coreSql

import (
	"github.com/jmoiron/sqlx"
	"github.com/konsultin/sqlk"
	"github.com/konsultin/sqlk/pq/query"
)

type User struct {
	GetUserByXid     *sqlx.Stmt
	GetUserById      *sqlx.Stmt
	FindByIdentifier *sqlx.Stmt
	Insert           *sqlx.NamedStmt
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
		GetUserById: db.MustPrepareRebind(
			query.Select(
				query.Column("*"),
			).
				From(UserSchema).
				Where(
					query.Equal(query.Column("id")),
				).Build(),
		),
		FindByIdentifier: db.MustPrepareRebind(`
			SELECT * FROM "user"
			WHERE email = ? OR phone = ? OR username = ?
			LIMIT 1
		`),
		Insert: db.MustPrepareNamed(`
			INSERT INTO "user" (
				xid, username, full_name, email, phone, age, avatar, status_id,
				created_at, updated_at, modified_by, version, metadata
			) VALUES (
				:xid, :username, :full_name, :email, :phone, :age, :avatar, :status_id,
				:created_at, :updated_at, :modified_by, :version, :metadata
			) RETURNING id
		`),
	}
}
