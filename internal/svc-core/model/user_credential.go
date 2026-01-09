package model

import (
	"database/sql"

	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/timek"
)

type UserCredential struct {
	Id               int64                 `db:"id"`
	UserId           int64                 `db:"user_id"`
	AuthProviderId   dto.AuthProvider_Enum `db:"auth_provider_id"`
	CredentialKey    string                `db:"credential_key"`    // email/phone/username for PASSWORD, provider_user_id for OAuth
	CredentialSecret sql.NullString        `db:"credential_secret"` // password_hash for PASSWORD, null for OAuth
	IsVerified       bool                  `db:"is_verified"`
	VerifiedAt       sql.NullTime          `db:"verified_at"`
	CreatedAt        timek.Time            `db:"created_at"`
	UpdatedAt        timek.Time            `db:"updated_at"`

	User *User `db:"-"`
}

func NewUserCredential() *UserCredential {
	return &UserCredential{}
}
