package config

type Config struct {
	Port  uint16 `envconfig:"PORT" default:"8080"`
	Debug bool   `envconfig:"DEBUG" default:"false"`

	HttpClientRequestTimeout int `envconfig:"HTTP_CLIENT_REQUEST_TIMEOUT" default:"10"`

	DatabaseDriver          string `envconfig:"DB_DRIVER" default:"mysql"`
	DatabaseHost            string `envconfig:"DB_HOST" default:"localhost"`
	DatabasePort            string `envconfig:"DB_PORT" default:"3306"`
	DatabaseUsername        string `envconfig:"DB_USERNAME" default:"root"`
	DatabasePassword        string `envconfig:"DB_PASSWORD" default:""`
	DatabaseName            string `envconfig:"DB_NAME" default:""`
	DatabaseMaxIdleConn     int    `envconfig:"DB_MAX_IDLE_CONN" default:"10"`
	DatabaseMaxOpenConn     int    `envconfig:"DB_MAX_OPEN_CONN" default:"100"`
	DatabaseMaxConnLifetime int    `envconfig:"DB_MAX_CONN_LIFETIME" default:"300"`
}
