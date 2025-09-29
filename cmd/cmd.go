package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal"
	"github.com/chains-lab/sso-svc/internal/data"
	"github.com/chains-lab/sso-svc/internal/data/pgdb"
	"github.com/chains-lab/sso-svc/internal/domain"
	"github.com/chains-lab/sso-svc/internal/domain/services/session"
	"github.com/chains-lab/sso-svc/internal/domain/services/user"
	"github.com/chains-lab/sso-svc/internal/infra/jwtmanager"
	"github.com/chains-lab/sso-svc/internal/rest"
	"github.com/chains-lab/sso-svc/internal/rest/controller"
)

func StartServices(ctx context.Context, cfg internal.Config, log logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := data.NewDatabase(
		pgdb.NewUsers(pg),
		pgdb.NewSessions(pg),
	)

	jwtTokenManager := jwtmanager.NewManager(jwtmanager.Config{
		AccessSK:   cfg.JWT.User.AccessToken.SecretKey,
		RefreshSK:  cfg.JWT.User.RefreshToken.SecretKey,
		AccessTTL:  cfg.JWT.User.AccessToken.TokenLifetime,
		RefreshTTL: cfg.JWT.User.RefreshToken.TokenLifetime,
		Iss:        cfg.Service.Name,
	})

	logic := domain.NewCore(
		user.New(database),
		session.New(database, jwtTokenManager),
	)

	ctrl := controller.NewService(log, cfg.GoogleOAuth(), logic)

	run(func() { rest.Run(ctx, cfg, log, ctrl) })
}
