package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env   string `envconfig:"APP_ENV" default:"development"`
	Port  uint16 `envconfig:"PORT" default:"8080"`
	Debug bool   `envconfig:"DEBUG" default:"false"`

	HTTPReadTimeoutSeconds  int `envconfig:"HTTP_READ_TIMEOUT_SECONDS" default:"15"`
	HTTPWriteTimeoutSeconds int `envconfig:"HTTP_WRITE_TIMEOUT_SECONDS" default:"15"`
	HTTPIdleTimeoutSeconds  int `envconfig:"HTTP_IDLE_TIMEOUT_SECONDS" default:"60"`

	RateLimitRPS   int `envconfig:"RATE_LIMIT_RPS" default:"25"`
	RateLimitBurst int `envconfig:"RATE_LIMIT_BURST" default:"50"`

	CORSAllowOrigins []string `envconfig:"CORS_ALLOW_ORIGINS" default:"*"`

	DatabaseDriver          string `envconfig:"DB_DRIVER" default:"mysql"`
	DatabaseHost            string `envconfig:"DB_HOST" default:"localhost"`
	DatabasePort            string `envconfig:"DB_PORT" default:"3306"`
	DatabaseUsername        string `envconfig:"DB_USERNAME" default:"root"`
	DatabasePassword        string `envconfig:"DB_PASSWORD" default:""`
	DatabaseName            string `envconfig:"DB_NAME" default:""`
	DatabaseMaxIdleConn     int    `envconfig:"DB_MAX_IDLE_CONN" default:"10"`
	DatabaseMaxOpenConn     int    `envconfig:"DB_MAX_OPEN_CONN" default:"100"`
	DatabaseMaxConnLifetime int    `envconfig:"DB_MAX_CONN_LIFETIME" default:"300"`
	DatabaseTimeoutSeconds  int    `envconfig:"DB_TIMEOUT_SECONDS" default:"5"`
}

// Load reads environment variables (optionally from .env) into Config with defaults, and validates them.
func Load() (*Config, error) {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "" {
		env = strings.ToLower(os.Getenv("GO_ENV"))
	}
	if env == "" {
		env = "development"
	}

	// Only load .env in non-production to prevent accidental prod secrets leakage.
	if env != "production" {
		_ = godotenv.Load()
	}

	cfg := &Config{Env: env}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	cfg.Env = env

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Port == 0 || c.Port > 65535 {
		return fmt.Errorf("invalid PORT %d", c.Port)
	}

	if c.HTTPReadTimeoutSeconds <= 0 || c.HTTPWriteTimeoutSeconds <= 0 || c.HTTPIdleTimeoutSeconds <= 0 {
		return fmt.Errorf("HTTP timeouts must be greater than zero")
	}

	if c.RateLimitRPS <= 0 || c.RateLimitBurst <= 0 {
		return fmt.Errorf("rate limit values must be greater than zero")
	}

	driver := strings.ToLower(c.DatabaseDriver)
	switch driver {
	case "mysql", "mariadb":
		c.DatabaseDriver = "mysql"
	case "postgres", "postgresql", "pg":
		c.DatabaseDriver = "postgres"
	default:
		return fmt.Errorf("unsupported DB_DRIVER '%s'", c.DatabaseDriver)
	}

	if c.DatabaseHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DatabaseUsername == "" {
		return fmt.Errorf("DB_USERNAME is required")
	}
	if c.DatabaseName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.DatabaseTimeoutSeconds <= 0 {
		return fmt.Errorf("DB_TIMEOUT_SECONDS must be greater than zero")
	}

	return nil
}
