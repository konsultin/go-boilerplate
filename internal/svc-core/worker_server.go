package svcCore

import (
	"context"

	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/project-goes-here/internal/svc-core/constant"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/project-goes-here/internal/svc-core/service"
	"github.com/konsultin/logk"
	logkOption "github.com/konsultin/logk/option"
	"github.com/nats-io/nats.go"
)

func (s *Server) InitWorker() {
	// Subscribing to example event
	s.nats.Subscribe(constant.JobExample, s.HandleExampleWorker)

	s.log.Info("Worker initialized and listening...")
}

func (s *Server) HandleExampleWorker(msg *nats.Msg) {
	s.log.Infof("[WORKER] Received message from Repo: %s", string(msg.Data))
}

// NewWorkerService creates a service instance for worker context.
// It is similar to NewService but adapted for non-HTTP contexts (no fasthttp.RequestCtx).
func (s *Server) NewWorkerService(ctx context.Context) (*service.Service, error) {
	rc, err := s.repo.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return s.svc.
		WithContext(ctx).
		WithRepo(rc).
		WithLog(logk.Get().NewChild(logkOption.WithNamespace(constant.ServiceName+"/worker"), logkOption.Context(ctx))).
		WithSubject(&model.Subject{
			Id:       "SYSTEM",
			FullName: "System-Worker",
			Role:     dto.Role_Enum_name[int32(dto.Role_ADMIN)],
		}), nil
}
