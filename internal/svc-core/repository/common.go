package repository

import (
	"github.com/Konsultin/project-goes-here/config"
	"github.com/Konsultin/project-goes-here/libs/errk"
	"github.com/Konsultin/project-goes-here/libs/logk"
	logkOption "github.com/Konsultin/project-goes-here/libs/logk/option"
	"github.com/Konsultin/project-goes-here/libs/sqlk"
)

type Repository struct {
	config   *RepositoryConfig
	db       *sqlk.Database
	adapters *repositoryAdapters
}

func NewRepository(cfg *config.Config) (*Repository, error) {

	db, err := sqlk.NewDatabase(sqlk.Config{
		Driver:          cfg.DatabaseDriver,
		Host:            cfg.DatabaseHost,
		Port:            cfg.DatabasePort,
		Username:        cfg.DatabaseUsername,
		Password:        cfg.DatabasePassword,
		Database:        cfg.DatabaseName,
		MaxIdleConn:     &cfg.DatabaseMaxIdleConn,
		MaxOpenConn:     &cfg.DatabaseMaxOpenConn,
		MaxConnLifetime: &cfg.DatabaseMaxConnLifetime,
	})

	if err != nil {
		logk.Get().Error("Failed to connect to database", logkOption.Error(err))
		return nil, errk.Trace(err)
	}

	// Init repository config
	repoConfig, err := NewRepositoryConfig(cfg)
	if err != nil {
		logk.Get().Error("Failed to initialize repository config", logkOption.Error(err))
		return nil, errk.Trace(err)
	}

	adapters, err := newRepositoryAdapters(cfg)
	if err != nil {
		logk.Get().Error("Failed to initialize repository adapters", logkOption.Error(err))
		return nil, errk.Trace(err)
	}

	var r = Repository{
		config:   repoConfig,
		db:       db,
		adapters: adapters,
	}

	return &r, nil
}
