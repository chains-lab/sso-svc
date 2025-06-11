package cli

import (
	"context"
	"sync"

	"github.com/chains-lab/chains-auth/internal/app"
	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/chains-lab/chains-auth/internal/rest"
	"github.com/sirupsen/logrus"
)

func Start(ctx context.Context, cfg config.Config, log *logrus.Logger, wg *sync.WaitGroup, app *app.App) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	restApi := rest.NewRest(cfg, log, app)
	run(func() { restApi.Run(ctx, log) })
}
