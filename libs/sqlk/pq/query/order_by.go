package query

import (
	"github.com/Konsultin/project-goes-here/libs/sqlk"
	"github.com/Konsultin/project-goes-here/libs/sqlk/op"
)

type orderByWriter struct {
	sqlk.ColumnWriter
	direction op.SortDirection
}

func (o *orderByWriter) OrderByQuery() string {
	var direction string
	if o.direction == op.Descending {
		direction = "DESC"
	} else {
		direction = "ASC"
	}
	return o.ColumnQuery() + " " + direction
}
