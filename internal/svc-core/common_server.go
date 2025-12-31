package svcCore

import (
	"time"

	"github.com/konsultin/project-goes-here/config"
	"github.com/konsultin/project-goes-here/internal/svc-core/constant"
	"github.com/konsultin/project-goes-here/internal/svc-core/repository"
	"github.com/konsultin/project-goes-here/internal/svc-core/service"
	"github.com/konsultin/project-goes-here/libs/errk"
	"github.com/konsultin/project-goes-here/libs/logk"
	logkOption "github.com/konsultin/project-goes-here/libs/logk/option"
)

type Server struct {
	config    *config.Config
	startedAt time.Time
	svc       *service.Service
	repo      *repository.Repository
	log       logk.Logger
}

func New(config *config.Config, startedAt time.Time) (*Server, error) {
	repo, err := repository.NewRepository(config)
	if err != nil {
		return nil, errk.Trace(err)
	}

	svc := service.NewService(repo)

	server := &Server{
		config:    config,
		startedAt: startedAt,
		svc:       svc,
		repo:      repo,
		log:       logk.Get().NewChild(logkOption.WithNamespace(constant.ServiceName + "/server")),
	}

	return server, nil

}

func (s *Server) Close() error {
	return s.repo.Close()
}
