package sqlk

const (
	DriverMySQL      = "mysql"
	DriverPostgreSQL = "postgres"

	Null = `null`
)

const (
	DefaultMaxIdleConn     = 10
	DefaultMaxOpenConn     = 10
	DefaultMaxConnLifetime = 1
)

const (
	Separator = ", "
)
