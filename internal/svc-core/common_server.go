package svcCore

import (
	"fmt"
	"time"

	"github.com/Konsultin/project-goes-here/config"
	"github.com/Konsultin/project-goes-here/internal/svc-core/constant"
	"github.com/Konsultin/project-goes-here/internal/svc-core/repository"
	"github.com/Konsultin/project-goes-here/internal/svc-core/service"
	"github.com/Konsultin/project-goes-here/libs/logk"
	logkOption "github.com/Konsultin/project-goes-here/libs/logk/option"
	f "github.com/valyala/fasthttp"
)

type Server struct {
	config    *config.Config
	startedAt time.Time
	svc       *service.Service
	repo      *repository.Repository
	log       logk.Logger
}

func New(config *config.Config, startedAt time.Time) *Server {
	repo, err := repository.NewRepository(config)
	if err != nil {
		panic(err)
	}

	svc := service.NewService(repo)

	server := &Server{
		config:    config,
		startedAt: startedAt,
		svc:       svc,
		repo:      repo,
		log:       logk.Get().NewChild(logkOption.WithNamespace(constant.ServiceName + "/server")),
	}

	return server

}

func (s *Server) HealthCheck(ctx *f.RequestCtx) {
	uptime := time.Since(s.startedAt).String()
	response := fmt.Sprintf("Konsultin API is running. Uptime: %s", uptime)

	s.log.Debugf("Ran Health Check: %+s", response)

	ctx.SetStatusCode(f.StatusOK)
	ctx.SetBodyString(response)
}
