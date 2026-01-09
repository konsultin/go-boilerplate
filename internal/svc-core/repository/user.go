package repository

import (
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/errk"
)

func (r *Repository) FindUserByXid(xid string) (*model.User, error) {
	var m model.User
	err := r.sql.User.GetUserByXid.GetContext(r.ctx, &m, xid)
	if err != nil {
		return nil, errk.Trace(err)
	}
	return &m, nil
}

func (r *Repository) FindUserById(id int64) (*model.User, error) {
	var m model.User
	err := r.sql.User.GetUserById.GetContext(r.ctx, &m, id)
	if err != nil {
		return nil, errk.Trace(err)
	}
	return &m, nil
}

func (r *Repository) FindUserByIdentifier(identifier string) (*model.User, error) {
	var m model.User
	// Pass identifier 3 times for email, phone, username comparison
	err := r.sql.User.FindByIdentifier.GetContext(r.ctx, &m, identifier, identifier, identifier)
	if err != nil {
		return nil, errk.Trace(err)
	}
	return &m, nil
}

func (r *Repository) InsertUser(user *model.User) error {
	err := r.sql.User.Insert.GetContext(r.ctx, &user.Id, user)
	if err != nil {
		return errk.Trace(err)
	}
	return nil
}
