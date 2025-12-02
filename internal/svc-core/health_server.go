package svcCore

import (
	"context"
	"time"

	"github.com/Konsultin/project-goes-here/dto"
	f "github.com/valyala/fasthttp"
)

func (s *Server) HealthCheck(ctx *f.RequestCtx) {
	checkCtx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.DatabaseTimeoutSeconds)*time.Second)
	defer cancel()

	deps := map[string]string{
		"database": "ok",
	}

	statusCode := f.StatusOK
	statusText := "ok"
	code := dto.CodeOK
	if err := s.repo.Ping(checkCtx); err != nil {
		statusCode = f.StatusServiceUnavailable
		statusText = "degraded"
		code = dto.CodeServiceUnavailable
		deps["database"] = err.Error()
		s.log.Warnf("health check failed: %v", err)
	}

	uptime := time.Since(s.startedAt).String()
	message := "Health check passed"
	if statusCode != f.StatusOK {
		message = "dependency check failed"
	}

	if statusCode != f.StatusOK && !s.config.Debug {
		deps["database"] = "unavailable"
	}

	data := dto.HealthData{
		Status:       statusText,
		Uptime:       uptime,
		Started:      s.startedAt.UTC().Format(time.RFC3339),
		Env:          s.config.Env,
		Hostname:     string(ctx.Request.URI().Host()),
		Dependencies: deps,
	}

	resp := dto.Response[dto.HealthData]{
		Message:   message,
		Code:      code,
		Data:      data,
		Timestamp: time.Now().UTC().UnixMilli(),
	}

	s.log.Debugf("Ran Health Check: %+v", resp)

	s.response(ctx, statusCode, resp)
}
