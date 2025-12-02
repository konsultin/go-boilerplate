package sqlk

import "fmt"

type Config struct {
	Driver          string
	Host            string
	Port            string
	Username        string
	Password        string
	Database        string
	MaxIdleConn     *int
	MaxOpenConn     *int
	MaxConnLifetime *int
}

func (c *Config) normalizeValue() {
	// Check for optional values, set values if unset
	if c.MaxIdleConn == nil {
		c.MaxIdleConn = NewInt(DefaultMaxIdleConn)
	}
	if c.MaxOpenConn == nil {
		c.MaxOpenConn = NewInt(DefaultMaxOpenConn)
	}
	if c.MaxConnLifetime == nil {
		c.MaxConnLifetime = NewInt(DefaultMaxConnLifetime)
	}

	// Normalize driver
	switch c.Driver {
	case "postgresql", "pg":
		c.Driver = DriverPostgreSQL
	case "mysql", "mariadb":
		c.Driver = DriverMySQL
	}
}
func (c *Config) getDSN() (dsn string, err error) {
	switch c.Driver {
	case DriverMySQL:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.Username, c.Password, c.Host, c.Port,
			c.Database)
	case DriverPostgreSQL:
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port,
			c.Username, c.Password, c.Database)
	default:
		err = fmt.Errorf("sqlk: unsupported database driver '%s'", c.Driver)
	}
	return
}
