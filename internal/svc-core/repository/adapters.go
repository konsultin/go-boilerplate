package repository

import (
	"time"

	"github.com/konsultin/project-goes-here/config"
	"github.com/konsultin/project-goes-here/internal/svc-core/constant"
	coreSql "github.com/konsultin/project-goes-here/internal/svc-core/sql"
)

type repositoryAdapters struct {
	jakartaLoc *time.Location
	sql        *coreSql.Statements
}

func newRepositoryAdapters(_ *config.Config) (*repositoryAdapters, error) {
	a := new(repositoryAdapters)

	loc, err := time.LoadLocation(constant.JakartaLocale)
	if err != nil {
		return nil, err
	}

	a.jakartaLoc = loc

	return a, nil
}
