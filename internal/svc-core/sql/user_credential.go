package coreSql

import (
	"github.com/jmoiron/sqlx"
	"github.com/konsultin/sqlk"
	"github.com/konsultin/sqlk/pq/query"
)

type UserCredentialSql struct {
	FindByProviderAndKey *sqlx.Stmt
	FindByUserId         *sqlx.Stmt
	Insert               *sqlx.NamedStmt
	UpdateSecret         *sqlx.Stmt
}

func NewUserCredential(db *sqlk.DatabaseContext) *UserCredentialSql {
	return &UserCredentialSql{
		FindByProviderAndKey: db.MustPrepareRebind(
			query.Select(
				query.Column("*"),
			).
				From(UserCredentialSchema).
				Where(
					query.And(
						query.Equal(query.Column("auth_provider_id")),
						query.Equal(query.Column("credential_key")),
					),
				).Build(),
		),
		FindByUserId: db.MustPrepareRebind(
			query.Select(
				query.Column("*"),
			).
				From(UserCredentialSchema).
				Where(
					query.Equal(query.Column("user_id")),
				).Build(),
		),
		Insert: db.MustPrepareNamed(`
			INSERT INTO user_credential (
				user_id, auth_provider_id, credential_key, credential_secret,
				is_verified, verified_at, created_at, updated_at
			) VALUES (
				:user_id, :auth_provider_id, :credential_key, :credential_secret,
				:is_verified, :verified_at, :created_at, :updated_at
			) RETURNING id
		`),
		UpdateSecret: db.MustPrepareRebind(`
			UPDATE user_credential
			SET credential_secret = ?, updated_at = NOW()
			WHERE id = ?
		`),
	}
}
