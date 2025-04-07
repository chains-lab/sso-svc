package cli

import (
	"context"
	"sync"

	"github.com/hs-zavet/sso-oauth/internal/api"
	"github.com/hs-zavet/sso-oauth/internal/app"
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/sirupsen/logrus"
)

func runServices(ctx context.Context, wg *sync.WaitGroup, app app.App, cfg config.Config, log *logrus.Logger) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	API := api.NewAPI(cfg, log)
	run(func() { API.Run(ctx, &app) })

	//run(func() { eventlistener.NewListener(cfg, &app) })
}
