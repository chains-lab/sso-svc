package cli

import (
	"context"
	"sync"

	"github.com/hs-zavet/sso-oauth/internal/api"
	"github.com/hs-zavet/sso-oauth/internal/app"
	"github.com/hs-zavet/sso-oauth/internal/config"
)

func runServices(ctx context.Context, wg *sync.WaitGroup, app app.App, cfg *config.Config) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	API := api.NewAPI(cfg)
	run(func() { API.Run(ctx, &app) })

	//run(func() { eventlistener.NewListener(cfg, &app) })
}
