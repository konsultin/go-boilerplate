package sqlk

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	config *Config
	dsn    string
	conn   *sqlx.DB
}

func NewDatabase(config Config) (*Database, error) {
	// Set default connection values
	config.normalizeValue()

	// Generate DSN
	dsn, err := config.getDSN()
	if err != nil {
		return nil, err
	}

	// Set config
	db := Database{
		config: &config,
		dsn:    dsn,
	}
	return &db, nil
}

func (db *Database) InitContext(ctx context.Context) error {
	// Create connection
	var conn *sqlx.DB
	var err error

	conn, err = sqlx.ConnectContext(ctx, db.config.Driver, db.dsn)

	if err != nil {
		return err
	}

	conn.SetConnMaxLifetime(time.Duration(*db.config.MaxConnLifetime) * time.Second)
	conn.SetMaxOpenConns(*db.config.MaxOpenConn)
	conn.SetMaxIdleConns(*db.config.MaxIdleConn)

	db.conn = conn
	return nil
}

func (db *Database) Close() error {
	if db.conn == nil {
		return nil
	}
	return db.conn.Close()
}

func (db *Database) WithContext(ctx context.Context) *DatabaseContext {
	return &DatabaseContext{
		conn: db.conn,
		ctx:  ctx,
	}
}

func (db *Database) Init() error {
	// Create connection
	conn, err := sqlx.Connect(db.config.Driver, db.dsn)
	if err != nil {
		return err
	}

	// Set connection settings
	conn.SetConnMaxLifetime(time.Duration(*db.config.MaxConnLifetime) * time.Second)
	conn.SetMaxOpenConns(*db.config.MaxOpenConn)
	conn.SetMaxIdleConns(*db.config.MaxIdleConn)

	// Set connection
	db.conn = conn

	return nil
}

func (db *Database) IsConnected() bool {
	return db.conn != nil
}

func (db *Database) GetConnection(ctx context.Context) (*sqlx.Conn, error) {
	return db.conn.Connx(ctx)
}
