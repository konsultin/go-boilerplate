package repository

import (
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/project-goes-here/libs/errk"
)

func (r *Repository) FindUserByXid(xid string) (*model.User, error) {
	var m model.User
	err := r.sql.User.GetUserByXid.GetContext(r.ctx, &m, xid)
	if err != nil {
		return nil, errk.Trace(err)
	}
	return &m, nil
}
