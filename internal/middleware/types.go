package middleware

import (
	"github.com/konsultin/logk"
	"github.com/konsultin/routek"
	"github.com/valyala/fasthttp"
)

type ErrorResponder func(ctx *fasthttp.RequestCtx, status int, code routek.Code, message string, err error)

type Config struct {
	Handler          fasthttp.RequestHandler
	Logger           logk.Logger
	OnError          ErrorResponder
	RateLimitRPS     int
	RateLimitBurst   int
	CORSAllowOrigins []string
	Metrics          *Metrics
}
