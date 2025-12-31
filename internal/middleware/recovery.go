package middleware

import (
	"fmt"

	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/project-goes-here/libs/errk"
	"github.com/konsultin/project-goes-here/libs/logk"
	logkOption "github.com/konsultin/project-goes-here/libs/logk/option"
	"github.com/valyala/fasthttp"
)

func Recovery(log logk.Logger, onError ErrorResponder) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			defer func() {
				if r := recover(); r != nil {
					panicErr := fmt.Errorf("%v", r)
					log.Error("panic recovered", logkOption.Error(errk.Trace(panicErr)))
					onError(ctx, fasthttp.StatusInternalServerError, dto.CodeInternalError, "internal server error", panicErr)
				}
			}()
			next(ctx)
		}
	}
}
