package handlers

import (
	"context"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/rest/presenter"
	"github.com/chains-lab/chains-auth/internal/app"
	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type App interface {
	UpdateUserRole(ctx context.Context, userID uuid.UUID, role, initiatorRole roles.Role) *ape.Error

	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, *ape.Error)
	GetUserByEmail(ctx context.Context, email string) (models.User, *ape.Error)

	GetSession(ctx context.Context, sessionID uuid.UUID) (models.Session, *ape.Error)
	GetUserSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, *ape.Error)

	Login(ctx context.Context, email, client string) (models.Session, models.TokensPair, *ape.Error)
	Refresh(ctx context.Context, userID, sessionID uuid.UUID, client, token string) (models.Session, models.TokensPair, *ape.Error)
	Logout(ctx context.Context, sessionID uuid.UUID) *ape.Error

	TerminateSessionsByOwner(ctx context.Context, userID uuid.UUID) *ape.Error
	DeleteSessionByOwner(ctx context.Context, sessionID, initiatorSessionID uuid.UUID) *ape.Error
	TerminateSessionsByAdmin(ctx context.Context, userID uuid.UUID) *ape.Error
	DeleteSessionByAdmin(ctx context.Context, sessionID, initiatorID, initiatorSessionID uuid.UUID) *ape.Error
}

type Presenter interface {
	InvalidPointer(w http.ResponseWriter, requestID uuid.UUID, err error)
	InvalidToken(w http.ResponseWriter, requestID uuid.UUID, err error)
	InvalidParameter(w http.ResponseWriter, requestID uuid.UUID, err error, parameter string)
	InvalidQuery(w http.ResponseWriter, requestID uuid.UUID, query string, err error)
	MismatchIdentification(w http.ResponseWriter, requestID uuid.UUID, parameter, pointer string)
	AppError(w http.ResponseWriter, requestID uuid.UUID, appErr *ape.Error)
}

type Handlers struct {
	app       App
	presenter Presenter
	log       *logrus.Entry
	cfg       config.Config
	google    oauth2.Config
}

func NewHandlers(cfg config.Config, log *logrus.Entry, app *app.App) Handlers {
	return Handlers{
		app:       app,
		cfg:       cfg,
		google:    config.InitGoogleOAuth(cfg),
		presenter: presenter.NewPresenter(log),
		log:       log,
	}
}
