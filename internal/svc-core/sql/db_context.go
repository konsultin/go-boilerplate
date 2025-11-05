package coreSql

import (
	"github.com/Konsultin/project-goes-here/internal/svc-core/model"
	"github.com/Konsultin/project-goes-here/libs/sqlk/schema"
)

// User Schemas
var (
	UserSchema = schema.New(schema.FromModelRef(new(model.User)), schema.As("user"))
)
