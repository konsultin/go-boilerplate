package coreSql

import (
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/project-goes-here/libs/sqlk/schema"
)

// User Schemas
var (
	UserSchema          = schema.New(schema.FromModelRef(new(model.User)), schema.As("user"))
	ClientAuthSchema    = schema.New(schema.FromModelRef(new(model.ClientAuth)), schema.As("clientAuth"))
	RoleSchema          = schema.New(schema.FromModelRef(new(model.Role)), schema.As("role"))
	RolePrivilegeSchema = schema.New(schema.FromModelRef(new(model.RolePrivilege)), schema.As("rolePrivilege"))
	PrivilegeSchema     = schema.New(schema.FromModelRef(new(model.Privilege)), schema.As("privilege"))
	AuthSessionSchema   = schema.New(schema.FromModelRef(new(model.AuthSession)), schema.As("auth_session"))
)
