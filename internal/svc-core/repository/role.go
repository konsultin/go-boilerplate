package repository

import (
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	coreSql "github.com/konsultin/project-goes-here/internal/svc-core/sql"
	"github.com/konsultin/errk"
	"github.com/konsultin/sqlk/option"
	"github.com/konsultin/sqlk/pq/query"
)

func (r *Repository) FindRoleById(id int32) (*model.Role, error) {
	var role model.Role
	err := r.sql.Role.FindById.GetContext(r.ctx, &role, id)
	if err != nil {
		return nil, errk.Trace(err)
	}
	return &role, nil
}

func (r *Repository) FindRolePrivilegeByRoleId(roleId int32) ([]model.RolePrivilegeJoinRow, error) {
	b := query.From(coreSql.RolePrivilegeSchema)

	// filters
	b = b.Select(
		query.Column("*"),
		query.Column("*", option.Schema(coreSql.RoleSchema)),
		query.Column("*", option.Schema(coreSql.PrivilegeSchema)),
	).
		Join(coreSql.RoleSchema, query.Equal(query.Column("roleId"), query.On("id", option.Schema(coreSql.RoleSchema)))).
		Join(coreSql.PrivilegeSchema, query.Equal(query.Column("privilegeId"), query.On("id", option.Schema(coreSql.PrivilegeSchema)))).
		Where(query.Equal(query.Column("roleId")))

	dbCtx := r.db.WithContext(r.ctx)
	selectQuery := dbCtx.Rebind(b.Build())

	// Execute query list
	var rows []model.RolePrivilegeJoinRow
	err := dbCtx.SelectContext(r.ctx, &rows, selectQuery, roleId)
	if err != nil {
		return nil, errk.Trace(err)
	}

	return rows, nil
}
