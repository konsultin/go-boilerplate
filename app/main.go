package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Konsultin/project-goes-here/config"
	"github.com/Konsultin/project-goes-here/dto"
	svcCore "github.com/Konsultin/project-goes-here/internal/svc-core"
	"github.com/Konsultin/project-goes-here/libs/logk"
	logkOption "github.com/Konsultin/project-goes-here/libs/logk/option"
	"github.com/Konsultin/project-goes-here/libs/routek"
	"github.com/valyala/fasthttp"
)

type responder struct {
	debug bool
}

func newResponder(debug bool) responder {
	return responder{debug: debug}
}

func (r responder) success(ctx *fasthttp.RequestCtx, status int, code dto.Code, message string, data any) {
	r.write(ctx, status, code, message, data)
}

func (r responder) error(ctx *fasthttp.RequestCtx, status int, code dto.Code, message string, err error) {
	var data any
	if r.debug && err != nil {
		data = map[string]any{"error": err.Error()}
	}
	r.write(ctx, status, code, message, data)
}

func (r responder) write(ctx *fasthttp.RequestCtx, status int, code dto.Code, message string, data any) {
	resp := dto.Response[any]{
		Message:   message,
		Code:      code,
		Data:      data,
		Timestamp: time.Now().UTC().UnixMilli(),
	}

	body, err := json.Marshal(resp)
	if err != nil {
		log.Printf("failed to marshal response: %v", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(`{"message":"internal server error","code":"INTERNAL_ERROR","data":null,"timestamp":0}`)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(status)
	ctx.SetBody(body)
}

type rateLimiter struct {
	tokens chan struct{}
}

func newRateLimiter(rps, burst int) *rateLimiter {
	rl := &rateLimiter{
		tokens: make(chan struct{}, burst),
	}

	// Pre-fill burst tokens.
	for i := 0; i < burst; i++ {
		rl.tokens <- struct{}{}
	}

	interval := time.Second / time.Duration(rps)
	if interval <= 0 {
		interval = time.Millisecond
	}
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
				// Bucket full, drop token.
			}
		}
	}()

	return rl
}

func (r *rateLimiter) allow() bool {
	select {
	case <-r.tokens:
		return true
	default:
		return false
	}
}

func konsultinAscii() string {
	return `
'     __  _   ___   ____   _____ __ __  _     ______  ____  ____       ___      ___ __ __ 
'    |  |/ ] /   \ |    \ / ___/|  |  || |   |      ||    ||    \     |   \    /  _]  |  |
'    |  ' / |     ||  _  (   \_ |  |  || |   |      | |  | |  _  |    |    \  /  [_|  |  |
'    |    \ |  O  ||  |  |\__  ||  |  || |___|_|  |_| |  | |  |  |    |  D  ||    _]  |  |
'    |     ||     ||  |  |/  \ ||  :  ||     | |  |   |  | |  |  | __ |     ||   [_|  :  |
'    |  .  ||     ||  |  |\    ||     ||     | |  |   |  | |  |  ||  ||     ||     |\   / 
'    |__|\_| \___/ |__|__| \___| \__,_||_____| |__|  |____||__|__||__||_____||_____| \_/  
'      
'    Boilerplate created by Kenly Krisaguino - @kenly.krisaguino on Instagram
'	 Version: 1.0.0
'                                                                                         
`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		logk.Get().Fatalf("Failed to load config: %v", err)
	}
	startedAt := time.Now()
	rootLog := logk.Get().NewChild(logkOption.WithNamespace("api"))
	rootLog.Infof("API starting... env=%s", cfg.Env)

	fmt.Println(konsultinAscii())

	coreServer := svcCore.New(cfg, startedAt)
	defer func() {
		if err := coreServer.Close(); err != nil {
			rootLog.Errorf("Failed to close resources: %v", err)
		}
	}()

	rt, err := routek.NewRouter(routek.Config{
		Handlers: map[string]any{
			"core": coreServer,
		},
	})
	if err != nil {
		rootLog.Fatalf("Failed to init router: %v", err)
	}

	responder := newResponder(cfg.Debug)
	limiter := newRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	handler := chainMiddleware(rt.Handler,
		recoveryMiddleware(rootLog, responder),
		loggingMiddleware(rootLog),
		rateLimitMiddleware(limiter, rootLog, responder),
		corsMiddleware(cfg.CORSAllowOrigins),
	)

	server := &fasthttp.Server{
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.HTTPReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTPWriteTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(cfg.HTTPIdleTimeoutSeconds) * time.Second,
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	errCh := make(chan error, 1)

	go func() {
		rootLog.Infof("Listening on %s", addr)
		if err := server.ListenAndServe(addr); err != nil {
			errCh <- err
		}
		close(errCh)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-stop:
		rootLog.Infof("Received signal %s, shutting down", sig)
	case err := <-errCh:
		if err != nil {
			rootLog.Fatalf("Server error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.ShutdownWithContext(shutdownCtx); err != nil {
		rootLog.Errorf("Graceful shutdown failed: %v", err)
	} else {
		rootLog.Info("Server stopped gracefully")
	}
}

func chainMiddleware(final fasthttp.RequestHandler, middlewares ...func(fasthttp.RequestHandler) fasthttp.RequestHandler) fasthttp.RequestHandler {
	handler := final
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func recoveryMiddleware(log logk.Logger, res responder) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("panic recovered: %v", r)
					res.error(ctx, fasthttp.StatusInternalServerError, dto.CodeInternalError, "internal server error", fmt.Errorf("%v", r))
				}
			}()
			next(ctx)
		}
	}
}

func loggingMiddleware(log logk.Logger) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			start := time.Now()
			next(ctx)
			duration := time.Since(start)
			log.Infof("%s %s -> %d in %s", ctx.Method(), ctx.Path(), ctx.Response.StatusCode(), duration)
		}
	}
}

func rateLimitMiddleware(rl *rateLimiter, log logk.Logger, res responder) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			if !rl.allow() {
				log.Warn("request rejected: rate limit exceeded")
				res.error(ctx, fasthttp.StatusTooManyRequests, dto.CodeTooManyRequests, "too many requests", nil)
				return
			}
			next(ctx)
		}
	}
}

func corsMiddleware(allowedOrigins []string) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			origin := string(ctx.Request.Header.Peek("Origin"))

			if originAllowed(origin, allowedOrigins) {
				if origin == "" && len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
					ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
				} else {
					ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
					ctx.Response.Header.Set("Vary", "Origin")
				}
				ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
			}

			ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			ctx.Response.Header.Set("Access-Control-Allow-Headers", "Authorization,Content-Type")

			if ctx.IsOptions() {
				ctx.SetStatusCode(fasthttp.StatusNoContent)
				return
			}

			next(ctx)
		}
	}
}

func originAllowed(origin string, allowed []string) bool {
	if origin == "" {
		return true
	}
	for _, allowedOrigin := range allowed {
		if allowedOrigin == "*" || strings.EqualFold(allowedOrigin, origin) {
			return true
		}
	}
	return false
}
