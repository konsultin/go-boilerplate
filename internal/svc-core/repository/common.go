package repository

import (
	"context"
	"fmt"

	"github.com/konsultin/project-goes-here/config"
	"github.com/konsultin/errk"
	"github.com/konsultin/logk"
	logkOption "github.com/konsultin/logk/option"
	"github.com/konsultin/natsk"
	"github.com/konsultin/sqlk"
)

type Repository struct {
	config  *RepositoryConfig
	db      *sqlk.Database
	log     logk.Logger
	isClone bool
	ctx     context.Context
	nats    *natsk.Client
	*repositoryAdapters
}

func NewRepository(cfg *config.Config, nats *natsk.Client) (*Repository, error) {

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
		logk.Get().Error("Failed to initialize database config", logkOption.Error(errk.Trace(err)))
		return nil, errk.Trace(err)
	}

	if err := db.Init(); err != nil {
		logk.Get().Error("Failed to connect to database", logkOption.Error(errk.Trace(err)))
		return nil, errk.Trace(err)
	}

	// Init repository config
	repoConfig, err := NewRepositoryConfig(cfg)
	if err != nil {
		logk.Get().Error("Failed to initialize repository config", logkOption.Error(errk.Trace(err)))
		return nil, errk.Trace(err)
	}

	adapters, err := newRepositoryAdapters(cfg, db)
	if err != nil {
		logk.Get().Error("Failed to initialize repository adapters", logkOption.Error(errk.Trace(err)))
		return nil, errk.Trace(err)
	}

	var r = Repository{
		config:             repoConfig,
		db:                 db,
		repositoryAdapters: adapters,
		nats:               nats,
		log:                logk.Get().NewChild(logkOption.WithNamespace("svc-core/repository")),
	}

	logk.Get().Infof("Connected to database '%s' successfully", cfg.DatabaseName)

	return &r, nil
}

func (r *Repository) Close() error {
	if r == nil || r.db == nil {
		return nil
	}
	if r.isClone {
		return nil
	}
	return r.db.Close()
}

func (r *Repository) Connect(ctx context.Context) (*Repository, error) {
	newR := *r
	newR.isClone = true
	newR.ctx = ctx
	return &newR, nil
}

func (r *Repository) Ping(ctx context.Context) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("repository not initialized")
	}
	return r.db.PingContext(ctx)
}
