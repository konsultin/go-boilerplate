package repository

import (
	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/errk"
)

// FindCredentialByKey finds a user credential by auth provider and key
func (r *Repository) FindCredentialByKey(authProviderId dto.AuthProvider_Enum, credentialKey string) (*model.UserCredential, error) {
	var credential model.UserCredential
	err := r.sql.UserCredential.FindByProviderAndKey.GetContext(r.ctx, &credential, authProviderId, credentialKey)
	if err != nil {
		return nil, errk.Trace(err)
	}
	return &credential, nil
}

// FindCredentialsByUserId finds all credentials for a user
func (r *Repository) FindCredentialsByUserId(userId int64) ([]*model.UserCredential, error) {
	var credentials []*model.UserCredential
	err := r.sql.UserCredential.FindByUserId.SelectContext(r.ctx, &credentials, userId)
	if err != nil {
		return nil, errk.Trace(err)
	}
	return credentials, nil
}

// InsertUserCredential inserts a new user credential
func (r *Repository) InsertUserCredential(credential *model.UserCredential) error {
	err := r.sql.UserCredential.Insert.GetContext(r.ctx, &credential.Id, credential)
	if err != nil {
		return errk.Trace(err)
	}
	return nil
}

// UpdateCredentialSecret updates the password hash for a credential
func (r *Repository) UpdateCredentialSecret(id int64, newSecret string) error {
	_, err := r.sql.UserCredential.UpdateSecret.ExecContext(r.ctx, newSecret, id)
	if err != nil {
		return errk.Trace(err)
	}
	return nil
}

// FindUserWithCredential finds user and their password credential by identifier
func (r *Repository) FindUserWithCredential(identifier string) (*model.User, *model.UserCredential, error) {
	user, err := r.FindUserByIdentifier(identifier)
	if err != nil {
		return nil, nil, errk.Trace(err)
	}

	// Find password credential for this user
	credential, err := r.FindCredentialByKey(dto.AuthProvider_PASSWORD, identifier)
	if err != nil {
		return nil, nil, errk.Trace(err)
	}

	return user, credential, nil
}
