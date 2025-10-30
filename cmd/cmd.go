package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal"
	"github.com/chains-lab/sso-svc/internal/data"
	"github.com/chains-lab/sso-svc/internal/domain/services/auth"
	"github.com/chains-lab/sso-svc/internal/domain/services/session"
	"github.com/chains-lab/sso-svc/internal/domain/services/user"
	"github.com/chains-lab/sso-svc/internal/events/listener"
	"github.com/chains-lab/sso-svc/internal/events/listener/callback"
	"github.com/chains-lab/sso-svc/internal/jwtmanager"
	"github.com/chains-lab/sso-svc/internal/passmanager"
	"github.com/chains-lab/sso-svc/internal/rest"
	"github.com/chains-lab/sso-svc/internal/rest/controller"
	"github.com/chains-lab/sso-svc/internal/rest/middlewares"
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

	database := data.NewDatabase(pg)

	jwtTokenManager := jwtmanager.NewManager(jwtmanager.Config{
		AccessSK:   cfg.JWT.User.AccessToken.SecretKey,
		RefreshSK:  cfg.JWT.User.RefreshToken.SecretKey,
		AccessTTL:  cfg.JWT.User.AccessToken.TokenLifetime,
		RefreshTTL: cfg.JWT.User.RefreshToken.TokenLifetime,
		Iss:        cfg.Service.Name,
	})
	passManager := passmanager.New()

	//eventPublisher := publisher.New(cfg.Kafka.Broker)

	userSvc := user.New(database)
	sessionSvc := session.New(database, jwtTokenManager)
	authSvc := auth.New(database, jwtTokenManager, passManager)

	ctrl := controller.New(log, cfg.GoogleOAuth(), userSvc, sessionSvc, authSvc)
	mdlv := middlewares.New(log)

	eventCallbacks := callback.NewService(authSvc)

	run(func() { rest.Run(ctx, cfg, log, mdlv, ctrl) })

	run(func() { listener.Run(ctx, log, cfg.Kafka.Broker, eventCallbacks) })
}
