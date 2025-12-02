package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Konsultin/project-goes-here/config"
	svcCore "github.com/Konsultin/project-goes-here/internal/svc-core"
	"github.com/Konsultin/project-goes-here/libs/logk"
	"github.com/Konsultin/project-goes-here/libs/routek"
	"github.com/valyala/fasthttp"
)

func KonsultinAscii() string {
	return `
'     __  _   ___   ____   _____ __ __  _     ______  ____  ____       ___      ___ __ __ 
'    |  |/ ] /   \ |    \ / ___/|  |  || |   |      ||    ||    \     |   \    /  _]  |  |
'    |  ' / |     ||  _  (   \_ |  |  || |   |      | |  | |  _  |    |    \  /  [_|  |  |
'    |    \ |  O  ||  |  |\__  ||  |  || |___|_|  |_| |  | |  |  |    |  D  ||    _]  |  |
'    |     ||     ||  |  |/  \ ||  :  ||     | |  |   |  | |  |  | __ |     ||   [_|  :  |
'    |  .  ||     ||  |  |\    ||     ||     | |  |   |  | |  |  ||  ||     ||     |\   / 
'    |__|\_| \___/ |__|__| \___| \__,_||_____| |__|  |____||__|__||__||_____||_____| \_/  
'                                                                                         
`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		logk.Get().Fatalf("Failed to load config: %v", err)
	}
	startedAt := time.Now()
	logk.Get().Info("API starting... (c) Konsultin")

	fmt.Println(KonsultinAscii())

	coreServer := svcCore.New(cfg, startedAt)

	rt, err := routek.NewRouter(routek.Config{
		Handlers: map[string]any{
			"core": coreServer,
		},
	})
	if err != nil {
		logk.Get().Fatalf("Failed to init router: %v", err)
		log.Fatalf("init router: %v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	if err := fasthttp.ListenAndServe(addr, rt.Handler); err != nil {
		logk.Get().Fatalf("Failed to Start Server: %v", err)
	}
}
