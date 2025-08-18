package admin

import (
	adminPoroto "github.com/chains-lab/sso-proto/gen/go/svc/admin"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/config"
)

type Service struct {
	app *app.App
	cfg config.Config

	adminPoroto.UnimplementedAdminServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}
