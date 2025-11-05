package query

import (
	"fmt"

	"github.com/Konsultin/project-goes-here/libs/sqlk"
	"github.com/Konsultin/project-goes-here/libs/sqlk/op"
	"github.com/Konsultin/project-goes-here/libs/sqlk/option"
	"github.com/Konsultin/project-goes-here/libs/sqlk/schema"
)

// whereCompareWriter

func newWhereComparisonWriter(col sqlk.ColumnWriter, operator op.Operator, args []interface{}) *whereCompareWriter {
	opts := option.EvaluateOptions(args)

	// Get variable writer
	v := opts.GetVariable(option.VariableKey)
	if v == nil {
		// Set default variable writer, by operator
		switch operator {
		case op.Equal, op.NotEqual, op.GreaterThan, op.GreaterThanEqual, op.LessThan, op.LessThanEqual, op.Like,
			op.NotLike, op.ILike, op.NotILike:
			v = new(bindVar)
		case op.Between, op.NotBetween:
			v = new(betweenBindVar)
		case op.Is, op.IsNot:
			v = new(nullVar)
		}
	}

	// Get alias
	as, _ := opts.GetString(option.AsKey)

	return &whereCompareWriter{
		ColumnWriter: col,
		op:           operator,
		variable:     v,
		as:           as,
	}
}

func newInWhereComparisonWriter(col sqlk.ColumnWriter, argCount int, operator op.Operator, args []interface{}) *whereCompareWriter {
	opts := option.EvaluateOptions(args)

	// Get variable writer
	v := opts.GetVariable(option.VariableKey)
	if v == nil {
		// Set default variable writer, by operator
		switch operator {
		case op.In, op.NotIn:
			v = &inBindVar{argCount: argCount}
		}
	}

	// Get alias
	as, _ := opts.GetString(option.AsKey)

	return &whereCompareWriter{
		ColumnWriter: col,
		op:           operator,
		variable:     v,
		as:           as,
	}
}

func Equal(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.Equal, args)
}

func NotEqual(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.NotEqual, args)
}

func GreaterThan(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.GreaterThan, args)
}

func GreaterThanEqual(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.GreaterThanEqual, args)
}

func LessThan(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.LessThan, args)
}

func LessThanEqual(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.LessThanEqual, args)
}

func Like(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.Like, args)
}

func NotLike(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.NotLike, args)
}

func ILike(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.ILike, args)
}

func NotILike(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.NotILike, args)
}

func Between(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.Between, args)
}

func NotBetween(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.NotBetween, args)
}

func In(col sqlk.ColumnWriter, argCount int, args ...interface{}) *whereCompareWriter {
	return newInWhereComparisonWriter(col, argCount, op.In, args)
}

func NotIn(col sqlk.ColumnWriter, argCount int, args ...interface{}) *whereCompareWriter {
	return newInWhereComparisonWriter(col, argCount, op.NotIn, args)
}

func IsNull(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.Is, args)
}

func IsNotNull(col sqlk.ColumnWriter, args ...interface{}) *whereCompareWriter {
	return newWhereComparisonWriter(col, op.IsNot, args)
}

// whereLogicWriter

func newWhereLogicalWriter(operator op.Operator, cn []sqlk.WhereWriter) *whereLogicWriter {
	return &whereLogicWriter{
		op:         operator,
		conditions: cn,
	}
}

func And(cn ...sqlk.WhereWriter) *whereLogicWriter {
	return newWhereLogicalWriter(op.And, cn)
}

func Or(cn ...sqlk.WhereWriter) *whereLogicWriter {
	return newWhereLogicalWriter(op.Or, cn)
}

func resolveFromTableFlag(ww sqlk.WhereWriter, from *schema.Schema) {
	// Switch type
	switch w := ww.(type) {
	case sqlk.WhereLogicWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			resolveFromTableFlag(cw, from)
		}
	case sqlk.WhereCompareWriter:
		// Get alias
		if w.GetTableName() == fromTableFlag {
			w.SetSchema(from)
		}
	}
}

func filterWhereWriters(ww sqlk.WhereWriter, tables map[schema.Reference]*schema.Schema) sqlk.WhereWriter {
	// Switch type
	switch w := ww.(type) {
	case sqlk.WhereLogicWriter:
		// Get conditions
		var conditions []sqlk.WhereWriter
		for _, cw := range w.GetConditions() {
			c := filterWhereWriters(cw, tables)
			// If no writer is set, then delete from array
			if c == nil {
				continue
			}
			conditions = append(conditions, c)
		}
		// Update conditions
		w.SetConditions(conditions)
	case sqlk.WhereCompareWriter:
		// Check if condition is registered in table
		table, ok := tables[w.GetSchemaRef()]
		if !ok {
			return nil
		}

		// Set alias
		w.SetTableAs(table.As())
	}
	return ww
}

func resolveJoinTableFlag(ww sqlk.WhereWriter, joinTable *schema.Schema) {
	// Switch type
	switch w := ww.(type) {
	case sqlk.WhereLogicWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			resolveJoinTableFlag(cw, joinTable)
		}
	case sqlk.WhereCompareWriter:
		// Get variable
		v := w.GetVariable()

		// Cast as column
		cv, ok := v.(sqlk.ColumnWriter)
		if !ok {
			return
		}

		// If join table flag reference is set, then set schema
		if cv.GetTableName() == joinTableFlag {
			cv.SetSchema(joinTable)
		}
	}
}

func setJoinTableAs(ww sqlk.WhereWriter, joinTable *schema.Schema, tableRefs map[schema.Reference]*schema.Schema) {
	switch w := ww.(type) {
	case sqlk.WhereLogicWriter:
		// Get conditions
		for _, cw := range w.GetConditions() {
			setJoinTableAs(cw, joinTable, tableRefs)
		}
	case sqlk.WhereCompareWriter:
		// Get column
		if cw, ok := w.(sqlk.ColumnWriter); ok {
			setJoinColumnTableAs(cw, joinTable, tableRefs)
		}

		// Get variable and set table alias
		v := w.GetVariable()
		if cv, ok := v.(sqlk.ColumnWriter); ok {
			setJoinColumnTableAs(cv, joinTable, tableRefs)
		}
	}
}

func setJoinColumnTableAs(column sqlk.ColumnWriter, joinTable *schema.Schema, tableRefs map[schema.Reference]*schema.Schema) {
	// Get arguments
	tableName := column.GetTableName()

	if tableName == skipTableFlag {
		return
	}

	col := column.GetColumn()

	// Check in joinSchema
	if tableName == joinTable.TableName() {
		if !joinTable.IsColumnExist(col) {
			panic(fmt.Errorf(`column "%s" is not declared in Table "%s"`, col, joinTable.TableName()))
		}
		// Set alias
		column.SetTableAs(joinTable.As())
		return
	}

	// Check in tableRefs
	sRef := column.GetSchemaRef()
	tRef, tRefOk := tableRefs[sRef]
	if !tRefOk {
		panic(fmt.Errorf(`table "%s" is not declared in Query Builder`, tableName))
	}
	// Check against column
	if !tRef.IsColumnExist(col) {
		panic(fmt.Errorf(`column "%s" is not declared in Table "%s"`, tableName, col))
	}
	column.SetTableAs(tRef.As())
}
