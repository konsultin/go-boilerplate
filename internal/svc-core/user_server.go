package svcCore

import (
	"github.com/konsultin/project-goes-here/dto"
	httpkPkg "github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk"
	unaryHttpk "github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk/unary"
	f "github.com/valyala/fasthttp"
)

func (s *Server) HandleCreateAnonymousUserSession(ctx *f.RequestCtx) (*dto.CreateUserSession_Result, error) {
	basicAuth := unaryHttpk.GetBasicAuth(ctx)
	if basicAuth == nil {
		s.log.Errorf("Basic Auth is not set in header")
		return nil, s.wrapError(ctx, httpkPkg.UnauthorizedError)
	}

	// Init Service
	svc, err := s.NewService(ctx)
	if err != nil {
		s.log.Errorf("Failed to create service: %v", err)
		return nil, err
	}

	// Call Service
	result, err := svc.CreateAnonymousUserSession(basicAuth, dto.Role_ANONYMOUS_USER)
	if err != nil {
		s.log.Errorf("Failed to create user session: %v", err)
		return nil, err
	}

	return result, nil
}

func (s *Server) HandleUserRefreshToken(ctx *f.RequestCtx) (*dto.CreateUserSession_Result_Data, error) {
	// Bind and validate request payload
	payload, err := httpkPkg.BindAndValidate[dto.RefreshSession_Payload](ctx)
	if err != nil {
		return nil, s.wrapError(ctx, err)
	}

	// Init Service
	svc, err := s.NewService(ctx)
	if err != nil {
		s.log.Errorf("Failed to create service: %v", err)
		return nil, err
	}
	defer svc.Close()

	data, err := svc.RefreshUserSession(payload)
	if err != nil {
		return nil, s.wrapError(ctx, err)
	}

	return data, nil
}
