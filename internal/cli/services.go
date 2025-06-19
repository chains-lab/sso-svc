package cli

import (
	"context"
	"sync"

	"github.com/chains-lab/chains-auth/internal/api"
	"github.com/chains-lab/chains-auth/internal/app"
	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/sirupsen/logrus"
)

func Start(ctx context.Context, cfg config.Config, log *logrus.Logger, wg *sync.WaitGroup, app *app.App) error {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := api.Run(ctx, cfg, log, app); err != nil {
			log.Fatalf("gRPC server exited with error: %v", err)
		}
	}()
	return nil
}
