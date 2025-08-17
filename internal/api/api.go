package api

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/api/grpc"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/logger"
	"golang.org/x/sync/errgroup"
)

func Start(ctx context.Context, cfg config.Config, log logger.Logger, app *app.App) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error { return grpc.Run(ctx, cfg, log, app) })

	return eg.Wait()
}
