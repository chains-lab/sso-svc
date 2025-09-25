package cmd

import (
	"context"
	"sync"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal"
	"github.com/chains-lab/sso-svc/internal/api"
	"github.com/chains-lab/sso-svc/internal/api/rest/controller"
	"github.com/chains-lab/sso-svc/internal/data"
	"github.com/chains-lab/sso-svc/internal/domain/services/session"
	"github.com/chains-lab/sso-svc/internal/domain/services/session/jwtmanager"
	"github.com/chains-lab/sso-svc/internal/domain/services/user"
)

func StartServices(ctx context.Context, cfg internal.Config, log logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	database := data.NewDatabase(cfg.Database.SQL.URL)
	sessionMod := session.NewService(database, jwtmanager.NewManager(cfg))
	userMod := user.NewService(database)

	Api := api.NewAPI(cfg, log)

	run(func() {
		handl := controller.NewService(cfg, log, userMod, sessionMod)

		Api.RunRest(ctx, handl)
	})
}
