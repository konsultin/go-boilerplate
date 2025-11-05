package query

import (
	"fmt"

	"github.com/Konsultin/project-goes-here/libs/sqlk"
	"github.com/Konsultin/project-goes-here/libs/sqlk/op"
	"github.com/Konsultin/project-goes-here/libs/sqlk/schema"
)

type joinWriter struct {
	method      op.JoinMethod
	table       *schema.Schema
	onCondition sqlk.WhereWriter
	index       int
}

func (j *joinWriter) GetSchemaRef() schema.Reference {
	return j.table.Ref()
}

func (j *joinWriter) GetIndex() int {
	return j.index
}

func (j *joinWriter) SetIndex(n int) {
	j.index = n
}

func (j *joinWriter) GetTableName() string {
	return j.table.TableName()
}

func (j *joinWriter) JoinQuery() string {
	// Get method query
	var method string
	switch j.method {
	case op.InnerJoin:
		method = "INNER JOIN"
	case op.RightJoin:
		method = "RIGHT JOIN"
	case op.FullJoin:
		method = "FULL OUTER JOIN"
	default:
		// Default to left join
		method = "LEFT JOIN"
	}

	// Generate table name
	table := j.table
	tableName := fmt.Sprintf(`"%s"`, table.TableName())
	if table.As() != "" {
		tableName += fmt.Sprintf(` AS "%s"`, table.As())
	}

	// Write condition
	condition := j.onCondition.WhereQuery()

	return fmt.Sprintf(`%s %s ON %s`, method, tableName, condition)
}
