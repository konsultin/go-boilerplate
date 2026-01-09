package repository

import (
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/errk"
)

func (r *Repository) FindClientAuthByClientId(id string) (*model.ClientAuth, error) {
	var m model.ClientAuth
	err := r.sql.ClientAuth.FindByClientId.GetContext(r.ctx, &m, id)
	if err != nil {
		return nil, errk.Trace(err)
	}
	return &m, nil
}
