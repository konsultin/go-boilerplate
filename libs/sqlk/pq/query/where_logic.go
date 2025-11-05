package query

import (
	"fmt"
	"strings"

	"github.com/Konsultin/project-goes-here/libs/sqlk"
	"github.com/Konsultin/project-goes-here/libs/sqlk/op"
)

type whereLogicWriter struct {
	op         op.Operator
	conditions []sqlk.WhereWriter
}

func (w *whereLogicWriter) SetConditions(conditions []sqlk.WhereWriter) {
	w.conditions = conditions
}

func (w *whereLogicWriter) WhereQuery() string {
	if len(w.conditions) == 0 {
		return ""
	}

	var separator string
	if w.op == op.Or {
		separator = " OR "
	} else {
		separator = " AND "
	}

	var conditions []string
	for _, cw := range w.conditions {
		// Create query
		cq := cw.WhereQuery()

		// If query
		if cq == "" {
			continue
		}

		// If condition is a logical, then add brackets
		if _, ok := cw.(sqlk.WhereLogicWriter); ok {
			cq = fmt.Sprintf("(%s)", cq)
		}

		conditions = append(conditions, cq)
	}

	return strings.Join(conditions, separator)
}

func (w *whereLogicWriter) GetConditions() []sqlk.WhereWriter {
	return w.conditions
}
