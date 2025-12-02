package repository

import (
	"time"

	"github.com/Konsultin/project-goes-here/config"
)

type RepositoryConfig struct {
	Timeout time.Duration
}

func NewRepositoryConfig(config *config.Config) (*RepositoryConfig, error) {
	repoConfig := new(RepositoryConfig)

	repoConfig.Timeout = time.Duration(config.DatabaseTimeoutSeconds) * time.Second

	return repoConfig, nil
}
