package svcCore

import (
	"context"
	"errors"
	"time"

	"github.com/konsultin/project-goes-here/config"
	"github.com/konsultin/project-goes-here/internal/svc-core/constant"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk"
	unaryHttpk "github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk/unary"
	"github.com/konsultin/project-goes-here/internal/svc-core/repository"
	"github.com/konsultin/project-goes-here/internal/svc-core/service"
	"github.com/konsultin/errk"
	"github.com/konsultin/logk"
	logkOption "github.com/konsultin/logk/option"
	"github.com/konsultin/natsk"
	"github.com/konsultin/sqlk"
	f "github.com/valyala/fasthttp"
)

type Server struct {
	config    *config.Config
	startedAt time.Time
	svc       *service.Service
	repo      *repository.Repository
	log       logk.Logger
	nats      *natsk.Client
}

func New(config *config.Config, startedAt time.Time) (*Server, error) {
	natsClient, err := natsk.New(config.NatsUrl)
	if err != nil {
		return nil, errk.Trace(err)
	}

	repo, err := repository.NewRepository(config, natsClient)
	if err != nil {
		return nil, errk.Trace(err)
	}

	svc := service.NewService(repo, config)

	server := &Server{
		config:    config,
		startedAt: startedAt,
		svc:       svc,
		repo:      repo,
		log:       logk.Get().NewChild(logkOption.WithNamespace(constant.ServiceName + "/server")),
		nats:      natsClient,
	}

	return server, nil

}

func (s *Server) Close() error {
	s.nats.Close()
	return s.repo.Close()
}

func (s *Server) NewService(ctx *f.RequestCtx) (*service.Service, error) {
	// Get subject from context
	subject := unaryHttpk.GetSubject(ctx)

	// Get db connection
	rc, err := s.repo.Connect(ctx)
	if err != nil {
		return nil, errk.Trace(err)
	}

	return s.svc.
		WithContext(ctx).
		WithRepo(rc).
		WithLog(logk.Get().NewChild(logkOption.WithNamespace(constant.ServiceName+"/service"), logkOption.Context(ctx))).
		WithSubject(&model.Subject{
			Id:       subject.Id,
			FullName: subject.FullName,
			Role:     subject.Role,
		}), nil
}

func (s *Server) wrapError(ctx *f.RequestCtx, err error) error {
	s.log.Errorf("Error returned from Service. ErrorType=%T %v", err, logkOption.Context(ctx), logkOption.Format(err))

	// Handle cancellation error
	if errors.Is(err, context.Canceled) {
		err = httpk.CancelError.Wrap(err)
	} else if sqlk.ErrorIsPqCancelStatementByUser(err) {
		err = httpk.CancelError.Wrap(err)
	}

	return err
}
