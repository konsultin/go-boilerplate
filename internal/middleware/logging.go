package middleware

import (
	"time"

	"github.com/konsultin/logk"
	"github.com/valyala/fasthttp"
)

func Logging(log logk.Logger, metrics *Metrics) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			start := time.Now()
			next(ctx)
			duration := time.Since(start)
			status := ctx.Response.StatusCode()

			if metrics != nil {
				metrics.Record(status, duration)
			}

			reqID := RequestIDFromContext(ctx)
			if reqID != "" {
				log.Infof("%s %s -> %d in %s req_id=%s", ctx.Method(), ctx.Path(), status, duration, reqID)
				return
			}

			log.Infof("%s %s -> %d in %s", ctx.Method(), ctx.Path(), status, duration)
		}
	}
}
