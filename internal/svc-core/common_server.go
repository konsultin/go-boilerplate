package svcCore

import (
	"context"
	"errors"
	"time"

	"github.com/go-konsultin/errk"
	"github.com/go-konsultin/logk"
	logkOption "github.com/go-konsultin/logk/option"
	"github.com/go-konsultin/natsk"
	"github.com/go-konsultin/sqlk"
	"github.com/konsultin/project-goes-here/config"
	"github.com/konsultin/project-goes-here/internal/svc-core/constant"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk"
	unaryHttpk "github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk/unary"
	"github.com/konsultin/project-goes-here/internal/svc-core/repository"
	"github.com/konsultin/project-goes-here/internal/svc-core/service"
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
	s.log.Errorf("Error returned from Service. ErrorType=%T Error=%+v", err, err)

	// Handle cancellation error
	if errors.Is(err, context.Canceled) {
		err = httpk.CancelError.Wrap(err)
	} else if sqlk.ErrorIsPqCancelStatementByUser(err) {
		err = httpk.CancelError.Wrap(err)
	}

	// Extract HTTP status from errk.Error metadata
	var errkErr *errk.Error
	if errors.As(err, &errkErr) {
		if status, ok := errkErr.Metadata()["http_status"].(int); ok {
			ctx.SetStatusCode(status)
		}
	}

	return err
}
