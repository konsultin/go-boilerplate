package svcCore

import (
	logkOption "github.com/konsultin/logk/option"
	f "github.com/valyala/fasthttp"
)

func (s *Server) HandleTriggerSimulation(ctx *f.RequestCtx) {
	// Create service
	svc, err := s.NewService(ctx)
	if err != nil {
		s.log.Error("Failed to create service", logkOption.Error(err))
		ctx.Error("Internal Server Error", f.StatusInternalServerError)
		return
	}
	defer svc.Close()

	// Run simulation
	if err := svc.RunSimulation(); err != nil {
		s.wrapError(ctx, err)
		return
	}

	ctx.SetStatusCode(f.StatusOK)
	ctx.SetBodyString("Simulation Triggered: Server -> Service -> Repo -> NATS -> Worker")
}
