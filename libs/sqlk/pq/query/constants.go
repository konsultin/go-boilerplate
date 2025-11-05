package query

const (
	forceWriteFlag = "__force__"     // Skip table reference checking
	fromTableFlag  = "__from__"      // Use table that is declared in from
	skipTableFlag  = "__skip__"      // Mark query part will be excluded
	joinTableFlag  = "__joinTable__" // Use join table that is declared in from

	AllColumns = "*"
)
