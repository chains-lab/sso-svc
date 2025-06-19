package handlers

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/api/presenter"
	"github.com/chains-lab/chains-auth/internal/app"
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type App interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)

	GetSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error)
	SelectUserSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, error)

	Login(ctx context.Context, email string, role roles.Role, client string) (models.Session, models.TokensPair, error)
	Refresh(ctx context.Context, userID, sessionID uuid.UUID, client, token string) (models.Session, models.TokensPair, error)

	DeleteSession(ctx context.Context, userID, sessionID uuid.UUID) error

	TerminateUserSessions(ctx context.Context, userID uuid.UUID) error

	AdminTerminateSessions(ctx context.Context, initiatorID, userID uuid.UUID) error
	AdminUpdateUserRole(ctx context.Context, initiatorID, userID uuid.UUID, role roles.Role) error
	AdminUpdateUserSubscription(ctx context.Context, initiatorID, userID, subscriptionID uuid.UUID) error
	AdminUpdateUserVerified(ctx context.Context, initiatorID, userID uuid.UUID, verified bool) error
	AdminUpdateUserSuspended(ctx context.Context, initiatorID, userID uuid.UUID, suspended bool) error
}

type Handlers struct {
	app       App
	log       *logrus.Entry
	cfg       config.Config
	google    oauth2.Config
	presenter presenter.Presenters

	sso.SsoServiceServer
}

func NewHandlers(cfg config.Config, log *logrus.Entry, app *app.App) Handlers {
	pres := presenter.NewPresenters(log)

	return Handlers{
		app:       app,
		cfg:       cfg,
		google:    config.InitGoogleOAuth(cfg),
		log:       log,
		presenter: pres,
	}
}

//func (h Handlers) mustEmbedUnimplementedSsoServiceServer() {
//	// This method is required to implement the SsoServiceServer interface.
//	// It can be left empty as it is not used in this context.
//}
