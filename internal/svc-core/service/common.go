package service

import (
	"github.com/Konsultin/project-goes-here/config"
	"github.com/Konsultin/project-goes-here/internal/svc-core/repository"
	"github.com/Konsultin/project-goes-here/libs/logk"
)

type Service struct {
	repo   *repository.Repository
	log    logk.Logger
	config *config.Config
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

