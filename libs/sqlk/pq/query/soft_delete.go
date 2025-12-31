package query

import (
	"fmt"

	"github.com/konsultin/project-goes-here/libs/sqlk"
	"github.com/konsultin/project-goes-here/libs/sqlk/op"
	"github.com/konsultin/project-goes-here/libs/sqlk/option"
	"github.com/konsultin/project-goes-here/libs/sqlk/schema"
)

// SoftDeleteBuilder builds UPDATE query for soft delete
type SoftDeleteBuilder struct {
	schema *schema.Schema
	where  sqlk.WhereWriter
}

// Build generates the UPDATE query that sets deleted_at
func (b *SoftDeleteBuilder) Build(args ...interface{}) string {
	// Get variable format option
	opts := option.EvaluateOptions(args)
	format, ok := opts.GetVariableFormat()
	if !ok {
		format = op.NamedVar
	}

	// Get soft delete column name
	deletedAtCol := b.schema.SoftDeleteColumn()

	// Set variable format in conditions
	if b.where == nil {
		// Set where to primary key
		b.where = Equal(Column(b.schema.PrimaryKey(), option.Schema(b.schema)))
	}

	// Set format in conditions
	setUpdateFormat(b.where, b.schema, format)

	// Write where
	where := b.where.WhereQuery()

	// Build UPDATE query that sets deleted_at = NOW()
	return fmt.Sprintf(`UPDATE "%s" SET "%s" = NOW() WHERE %s`, b.schema.TableName(), deletedAtCol, where)
}

// Where sets WHERE condition for soft delete
func (b *SoftDeleteBuilder) Where(w sqlk.WhereWriter) *SoftDeleteBuilder {
	b.where = w
	return b
}

// SoftDelete creates a soft delete builder for the given schema
func SoftDelete(s *schema.Schema) *SoftDeleteBuilder {
	if !s.SoftDelete() {
		panic(fmt.Errorf("soft delete is not enabled for schema %s", s.TableName()))
	}

	return &SoftDeleteBuilder{
		schema: s,
	}
}

// RestoreBuilder builds UPDATE query to restore soft-deleted records
type RestoreBuilder struct {
	schema *schema.Schema
	where  sqlk.WhereWriter
}

// Build generates the UPDATE query that sets deleted_at = NULL
func (b *RestoreBuilder) Build(args ...interface{}) string {
	// Get variable format option
	opts := option.EvaluateOptions(args)
	format, ok := opts.GetVariableFormat()
	if !ok {
		format = op.NamedVar
	}

	// Get soft delete column name
	deletedAtCol := b.schema.SoftDeleteColumn()

	// Set variable format in conditions
	if b.where == nil {
		// Set where to primary key
		b.where = Equal(Column(b.schema.PrimaryKey(), option.Schema(b.schema)))
	}

	// Set format in conditions
	setUpdateFormat(b.where, b.schema, format)

	// Write where
	where := b.where.WhereQuery()

	// Build UPDATE query that sets deleted_at = NULL
	return fmt.Sprintf(`UPDATE "%s" SET "%s" = NULL WHERE %s`, b.schema.TableName(), deletedAtCol, where)
}

// Where sets WHERE condition for restore
func (b *RestoreBuilder) Where(w sqlk.WhereWriter) *RestoreBuilder {
	b.where = w
	return b
}

// Restore creates a restore builder for the given schema
func Restore(s *schema.Schema) *RestoreBuilder {
	if !s.SoftDelete() {
		panic(fmt.Errorf("soft delete is not enabled for schema %s", s.TableName()))
	}

	return &RestoreBuilder{
		schema: s,
	}
}

// WithTrashed returns a WHERE condition that includes soft-deleted records (no filter)
func WithTrashed() sqlk.WhereWriter {
	// Return nil to indicate no additional filtering needed
	return nil
}

// OnlyTrashed returns a WHERE condition that only includes soft-deleted records
func OnlyTrashed(s *schema.Schema) sqlk.WhereWriter {
	if !s.SoftDelete() {
		return nil
	}

	deletedAtCol := s.SoftDeleteColumn()
	return IsNotNull(Column(deletedAtCol, option.Schema(s)))
}
