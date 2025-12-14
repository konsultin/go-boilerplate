package service

import (
	"github.com/Konsultin/project-goes-here/config"
	"github.com/Konsultin/project-goes-here/dto"
	"github.com/Konsultin/project-goes-here/internal/svc-core/model"
	"github.com/Konsultin/project-goes-here/internal/svc-core/repository"
	"github.com/Konsultin/project-goes-here/libs/logk"
)

type Service struct {
	repo    *repository.Repository
	log     logk.Logger
	config  *config.Config
	subject *model.Subject
}

func (s *Service) WithSubject(subject *dto.Subject) *Service {
	newS := *s
	newS.subject = model.NewSubject(subject)
	return &newS
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}
