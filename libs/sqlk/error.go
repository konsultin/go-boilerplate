package sqlk

import "strings"

var RowNotUpdatedError = new(rowNotUpdatedError)

type rowNotUpdatedError struct{}

func (n *rowNotUpdatedError) Error() string {
	return "sqlk: row is not updated"
}

func ErrorIsPqCancelStatementByUser(err error) bool {
	return strings.Contains(err.Error(), "pq: canceling statement due to user request")
}
