package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/konsultin/errk"
	"github.com/konsultin/logk"
	logkOption "github.com/konsultin/logk/option"
	"github.com/konsultin/project-goes-here/config"
	_ "github.com/konsultin/project-goes-here/docs" // Swagger docs
	"github.com/konsultin/project-goes-here/internal/middleware"
	svcCore "github.com/konsultin/project-goes-here/internal/svc-core"
	"github.com/konsultin/routek"
	fasthttpSwagger "github.com/swaggo/fasthttp-swagger"
	"github.com/valyala/fasthttp"
)

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
	 Version: 0.4.0
'                                                                                         
	`
}

// @title           Konsultin API Boilerplate
// @version         0.4.0
// @description     This is a sample server for Konsultin API Boilerplate.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		logk.Get().Fatal("Failed to load config", logkOption.Error(errk.Trace(err)))
		return
	}
	startedAt := time.Now()
	rootLog := logk.Get().NewChild(logkOption.WithNamespace("api"))
	rootLog.Infof("API starting... env=%s", cfg.Env)

	fmt.Println(konsultinAscii())

	coreServer, err := svcCore.New(cfg, startedAt)
	if err != nil {
		rootLog.Fatal("Failed to init core server", logkOption.Error(errk.Trace(err)))
		return
	}
	defer func() {
		if coreServer == nil {
			return
		}
		if err := coreServer.Close(); err != nil {
			rootLog.Error("Failed to close resources", logkOption.Error(errk.Trace(err)))
		}
	}()

	resp := routek.NewResponder(cfg.Debug)

	// Init Worker (NATS)
	coreServer.InitWorker()

	rt, err := routek.NewRouter(routek.Config{
		Handlers:  map[string]any{"core": coreServer},
		Responder: resp,
	})
	if err != nil {
		rootLog.Fatal("Failed to init router", logkOption.Error(errk.Trace(err)))
		return
	}

	// Swagger Handler Wrapper
	swaggerHandler := fasthttpSwagger.WrapHandler(fasthttpSwagger.URL("/swagger/doc.json"))

	handler, err := middleware.Init(middleware.Config{
		Handler: func(ctx *fasthttp.RequestCtx) {
			path := string(ctx.Path())
			if len(path) >= 8 && path[:8] == "/swagger" {
				swaggerHandler(ctx)
				return
			}
			rt.Handler(ctx)
		},
		Logger:           rootLog,
		OnError:          resp.Error,
		RateLimitRPS:     cfg.RateLimitRPS,
		RateLimitBurst:   cfg.RateLimitBurst,
		CORSAllowOrigins: cfg.CORSAllowOrigins,
	})
	if err != nil {
		rootLog.Fatal("Failed to init middleware", logkOption.Error(errk.Trace(err)))
		return
	}

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
			rootLog.Fatal("Server error", logkOption.Error(errk.Trace(err)))
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.ShutdownWithContext(shutdownCtx); err != nil {
		rootLog.Error("Graceful shutdown failed", logkOption.Error(errk.Trace(err)))
	} else {
		rootLog.Info("Server stopped gracefully")
	}
}
