package svcCore

import (
	"time"

	"github.com/konsultin/project-goes-here/dto"
	f "github.com/valyala/fasthttp"
)

func (s *Server) HealthCheck(ctx *f.RequestCtx) (*dto.HealthData, error) {
	uptime := time.Since(s.startedAt)

	data := dto.HealthData{
		Status:   "HEALTHY",
		Uptime:   uptime.String(),
		Started:  s.startedAt.UTC().Format(time.RFC3339),
		Env:      s.config.Env,
		Hostname: string(ctx.Request.URI().Host()),
	}

	s.log.Debugf("Ran Health Check: %+v", data)

	return &data, nil
}
