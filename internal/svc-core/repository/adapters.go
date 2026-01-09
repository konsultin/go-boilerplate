package repository

import (
	"context"
	"time"

	"github.com/konsultin/project-goes-here/config"
	"github.com/konsultin/project-goes-here/internal/svc-core/constant"
	coreSql "github.com/konsultin/project-goes-here/internal/svc-core/sql"
	"github.com/konsultin/sqlk"
)

type repositoryAdapters struct {
	jakartaLoc *time.Location
	sql        *coreSql.Statements
}

func newRepositoryAdapters(cfg *config.Config, db *sqlk.Database) (*repositoryAdapters, error) {
	a := new(repositoryAdapters)

	loc, err := time.LoadLocation(constant.JakartaLocale)
	if err != nil {
		return nil, err
	}

	a.jakartaLoc = loc
	a.sql = coreSql.New(db.WithContext(context.Background()))

	return a, nil
}
