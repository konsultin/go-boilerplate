package query

import (
	"fmt"
	"strings"

	"github.com/konsultin/project-goes-here/libs/sqlk"
	"github.com/konsultin/project-goes-here/libs/sqlk/schema"
)

// BulkInsertBuilder builds bulk INSERT queries
type BulkInsertBuilder struct {
	tableName string
	columns   []string
	values    []map[string]interface{}
	pk        string
}

// Values sets the values for bulk insert - accepts slice of maps
func (b *BulkInsertBuilder) Values(rows []map[string]interface{}) *BulkInsertBuilder {
	b.values = rows
	return b
}

// Build generates the bulk INSERT query
func (b *BulkInsertBuilder) Build() string {
	if len(b.values) == 0 {
		panic(fmt.Errorf("no values provided for bulk insert"))
	}

	count := len(b.columns)

	// Write columns
	columnQueries := make([]string, count)
	for i, col := range b.columns {
		columnQueries[i] = fmt.Sprintf(`"%s"`, col)
	}
	columns := strings.Join(columnQueries, sqlk.Separator)

	// Write values for each row
	rowQueries := make([]string, len(b.values))
	for i := range b.values {
		valueQueries := make([]string, count)
		for j, col := range b.columns {
			// Use named parameters with row index
			valueQueries[j] = fmt.Sprintf(`:row%d_%s`, i, col)
		}
		rowQueries[i] = fmt.Sprintf("(%s)", strings.Join(valueQueries, sqlk.Separator))
	}
	valuesClause := strings.Join(rowQueries, sqlk.Separator)

	// Compose returning
	returning := ""
	if b.pk != "" {
		returning = fmt.Sprintf(` RETURNING "%s"`, b.pk)
	}

	return fmt.Sprintf(`INSERT INTO "%s"(%s) VALUES %s%s`, b.tableName, columns, valuesClause, returning)
}

// BulkInsert creates a bulk insert builder
func BulkInsert(s *schema.Schema, column string, columnN ...string) *BulkInsertBuilder {
	// Init builder
	b := BulkInsertBuilder{
		tableName: s.TableName(),
		pk:        s.PrimaryKey(),
	}

	var columns []string
	if column == AllColumns {
		// Get all columns
		columns = s.InsertColumns()
	} else {
		inColumns := append([]string{column}, columnN...)
		for _, c := range inColumns {
			if s.IsColumnExist(c) {
				columns = append(columns, c)
			}
		}
	}

	// Set columns
	b.columns = columns

	return &b
}

// BulkUpdateBuilder builds bulk UPDATE queries using CASE statements
type BulkUpdateBuilder struct {
	schema      *schema.Schema
	columns     []string
	values      []map[string]interface{}
	primaryKeys []interface{}
}

// Values sets the values for bulk update - accepts slice of maps with primary key
// Each map must include the primary key field
func (b *BulkUpdateBuilder) Values(rows []map[string]interface{}) *BulkUpdateBuilder {
	b.values = rows
	pk := b.schema.PrimaryKey()

	// Extract primary keys
	b.primaryKeys = make([]interface{}, len(rows))
	for i, row := range rows {
		if pkVal, ok := row[pk]; ok {
			b.primaryKeys[i] = pkVal
		} else {
			panic(fmt.Errorf("primary key %s not found in row %d", pk, i))
		}
	}

	return b
}

// Build generates the bulk UPDATE query using CASE statements
func (b *BulkUpdateBuilder) Build() string {
	if len(b.values) == 0 {
		panic(fmt.Errorf("no values provided for bulk update"))
	}

	pk := b.schema.PrimaryKey()

	// Build CASE statements for each column
	caseStatements := make([]string, 0, len(b.columns))
	for _, col := range b.columns {
		if col == pk {
			continue // Skip primary key in SET clause
		}

		// Build CASE WHEN ... THEN ... END for this column
		whenClauses := make([]string, len(b.values))
		for i := range b.values {
			whenClauses[i] = fmt.Sprintf(`WHEN "%s" = :pk%d THEN :row%d_%s`, pk, i, i, col)
		}

		caseStmt := fmt.Sprintf(`"%s" = CASE %s END`, col, strings.Join(whenClauses, " "))
		caseStatements = append(caseStatements, caseStmt)
	}

	if len(caseStatements) == 0 {
		panic(fmt.Errorf("no columns to update in bulk update"))
	}

	setClause := strings.Join(caseStatements, sqlk.Separator)

	// Build WHERE clause with primary keys
	whereValues := make([]string, len(b.primaryKeys))
	for i := range b.primaryKeys {
		whereValues[i] = fmt.Sprintf(":pk%d", i)
	}
	whereClause := fmt.Sprintf(`"%s" IN (%s)`, pk, strings.Join(whereValues, sqlk.Separator))

	return fmt.Sprintf(`UPDATE "%s" SET %s WHERE %s`, b.schema.TableName(), setClause, whereClause)
}

// BulkUpdate creates a bulk update builder
func BulkUpdate(s *schema.Schema, column string, columnN ...string) *BulkUpdateBuilder {
	// Init builder
	b := BulkUpdateBuilder{
		schema: s,
	}

	var columns []string
	if column == AllColumns {
		// Get all columns
		columns = s.UpdateColumns()
	} else {
		inColumns := append([]string{column}, columnN...)
		pk := s.PrimaryKey()
		for _, c := range inColumns {
			if s.IsColumnExist(c) && c != pk {
				columns = append(columns, c)
			}
		}
	}

	// Set columns
	b.columns = columns

	return &b
}

// BulkDeleteBuilder builds bulk DELETE queries
type BulkDeleteBuilder struct {
	schema      *schema.Schema
	primaryKeys []interface{}
}

// IDs sets the primary key values for bulk delete
func (b *BulkDeleteBuilder) IDs(ids ...interface{}) *BulkDeleteBuilder {
	b.primaryKeys = ids
	return b
}

// Build generates the bulk DELETE query
func (b *BulkDeleteBuilder) Build() string {
	if len(b.primaryKeys) == 0 {
		panic(fmt.Errorf("no primary keys provided for bulk delete"))
	}

	pk := b.schema.PrimaryKey()

	// Build placeholders for IN clause
	placeholders := make([]string, len(b.primaryKeys))
	for i := range b.primaryKeys {
		placeholders[i] = fmt.Sprintf(":id%d", i)
	}

	whereClause := fmt.Sprintf(`"%s" IN (%s)`, pk, strings.Join(placeholders, sqlk.Separator))

	return fmt.Sprintf(`DELETE FROM "%s" WHERE %s`, b.schema.TableName(), whereClause)
}

// BulkDelete creates a bulk delete builder
func BulkDelete(s *schema.Schema) *BulkDeleteBuilder {
	// If soft delete is enabled, panic and direct user to use BulkSoftDelete() instead
	if s.SoftDelete() {
		panic(fmt.Errorf("soft delete is enabled for schema %s. Use query.BulkSoftDelete() instead", s.TableName()))
	}

	return &BulkDeleteBuilder{
		schema: s,
	}
}

// BulkSoftDeleteBuilder builds bulk soft DELETE queries
type BulkSoftDeleteBuilder struct {
	schema      *schema.Schema
	primaryKeys []interface{}
}

// IDs sets the primary key values for bulk soft delete
func (b *BulkSoftDeleteBuilder) IDs(ids ...interface{}) *BulkSoftDeleteBuilder {
	b.primaryKeys = ids
	return b
}

// Build generates the bulk soft DELETE query (UPDATE with deleted_at)
func (b *BulkSoftDeleteBuilder) Build() string {
	if len(b.primaryKeys) == 0 {
		panic(fmt.Errorf("no primary keys provided for bulk soft delete"))
	}

	pk := b.schema.PrimaryKey()
	deletedAtCol := b.schema.SoftDeleteColumn()

	// Build placeholders for IN clause
	placeholders := make([]string, len(b.primaryKeys))
	for i := range b.primaryKeys {
		placeholders[i] = fmt.Sprintf(":id%d", i)
	}

	whereClause := fmt.Sprintf(`"%s" IN (%s)`, pk, strings.Join(placeholders, sqlk.Separator))

	return fmt.Sprintf(`UPDATE "%s" SET "%s" = NOW() WHERE %s`, b.schema.TableName(), deletedAtCol, whereClause)
}

// BulkSoftDelete creates a bulk soft delete builder
func BulkSoftDelete(s *schema.Schema) *BulkSoftDeleteBuilder {
	if !s.SoftDelete() {
		panic(fmt.Errorf("soft delete is not enabled for schema %s", s.TableName()))
	}

	return &BulkSoftDeleteBuilder{
		schema: s,
	}
}
