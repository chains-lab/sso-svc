package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/chains-lab/kafkakit/box"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal"
	"github.com/chains-lab/sso-svc/internal/domain/modules/auth"
	"github.com/chains-lab/sso-svc/internal/events/producer"
	"github.com/chains-lab/sso-svc/internal/repo"
	"github.com/chains-lab/sso-svc/internal/rest"
	"github.com/chains-lab/sso-svc/internal/rest/controller"
	"github.com/chains-lab/sso-svc/internal/rest/middlewares"
	"github.com/chains-lab/sso-svc/internal/token"
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

	repository := repo.New(pg)

	kafkaBox := box.New(pg)

	jwtTokenManager := token.NewManager(token.Config{
		AccessSK:   cfg.JWT.User.AccessToken.SecretKey,
		RefreshSK:  cfg.JWT.User.RefreshToken.SecretKey,
		AccessTTL:  cfg.JWT.User.AccessToken.TokenLifetime,
		RefreshTTL: cfg.JWT.User.RefreshToken.TokenLifetime,
		Iss:        cfg.Service.Name,
	})

	kafkaProducer := producer.New(log, cfg.Kafka.Broker, kafkaBox)

	core := auth.NewService(repository, jwtTokenManager, kafkaProducer)

	ctrl := controller.New(log, cfg.GoogleOAuth(), core)
	mdlv := middlewares.New(log)

	run(func() { rest.Run(ctx, cfg, log, mdlv, ctrl) })

	run(func() { kafkaProducer.Run(ctx) })
}
