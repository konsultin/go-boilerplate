package query

import (
	"fmt"

	"github.com/Konsultin/project-goes-here/libs/sqlk"
	"github.com/Konsultin/project-goes-here/libs/sqlk/op"
	"github.com/Konsultin/project-goes-here/libs/sqlk/option"
	"github.com/Konsultin/project-goes-here/libs/sqlk/schema"
)

type DeleteBuilder struct {
	schema *schema.Schema
	where  sqlk.WhereWriter
}

func (b *DeleteBuilder) Build(args ...interface{}) string {
	// Get variable format option
	opts := option.EvaluateOptions(args)
	format, ok := opts.GetVariableFormat()
	if !ok {
		// If var format is not defined, then set default to query.NamedVar
		format = op.BindVar
	}

	// Set variable format in conditions
	if b.where == nil {
		// Set where to id
		b.where = Equal(Column(b.schema.PrimaryKey(), option.Schema(b.schema)))
	}

	// Set format in conditions
	setUpdateFormat(b.where, b.schema, format)

	// Write where
	where := b.where.WhereQuery()

	return fmt.Sprintf(`DELETE FROM "%s" WHERE %s`, b.schema.TableName(), where)
}

func (b *DeleteBuilder) Where(w sqlk.WhereWriter) *DeleteBuilder {
	b.where = w
	return b
}

func Delete(s *schema.Schema) *DeleteBuilder {
	return &DeleteBuilder{
		schema: s,
	}
}
