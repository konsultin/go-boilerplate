package middleware

import (
	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/project-goes-here/libs/logk"
	"github.com/valyala/fasthttp"
)

type ErrorResponder func(ctx *fasthttp.RequestCtx, status int, code dto.Code, message string, err error)

type Config struct {
	Handler          fasthttp.RequestHandler
	Logger           logk.Logger
	OnError          ErrorResponder
	RateLimitRPS     int
	RateLimitBurst   int
	CORSAllowOrigins []string
	Metrics          *Metrics
}
