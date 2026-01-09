package service

import (
	"context"

	"github.com/konsultin/project-goes-here/config"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/project-goes-here/internal/svc-core/repository"
	"github.com/konsultin/logk"
	logkOption "github.com/konsultin/logk/option"
)

type Service struct {
	repo    *repository.Repository
	log     logk.Logger
	ctx     context.Context
	config  *config.Config
	subject *model.Subject
}

func (s *Service) WithSubject(subject *model.Subject) *Service {
	newS := *s
	newS.subject = subject
	return &newS
}

func (s *Service) WithContext(ctx context.Context) *Service {
	newS := *s
	newS.ctx = ctx
	return &newS
}

func (s *Service) WithConfig(config *config.Config) *Service {
	newS := *s
	newS.config = config
	return &newS
}

func (s *Service) WithRepo(repo *repository.Repository) *Service {
	newS := *s
	newS.repo = repo
	return &newS
}

func (s *Service) WithLog(log logk.Logger) *Service {
	newS := *s
	newS.log = log
	return &newS
}

func NewService(repo *repository.Repository, config *config.Config) *Service {
	return &Service{
		repo:   repo,
		config: config,
	}
}

func (s *Service) Close() {
	// Returns connection to pool
	err := s.repo.Close()
	if err != nil {
		s.log.Error("Failed to close connection", logkOption.Error(err))
	} else {
		s.log.Tracef("DB: Connection returned to pool")
	}
}
