package handlers

import (
	"context"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/rest/controllers"
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
	CreateAccount(ctx context.Context, email string, role roles.Role) *ape.Error
	UpdateAccountRole(ctx context.Context, accountID uuid.UUID, role, initiatorRole roles.Role) *ape.Error

	GetAccountByID(ctx context.Context, accountID uuid.UUID) (models.Account, *ape.Error)
	GetAccountByEmail(ctx context.Context, email string) (models.Account, *ape.Error)

	GetSession(ctx context.Context, sessionID uuid.UUID) (models.Session, *ape.Error)
	GetSessions(ctx context.Context, accountID uuid.UUID) ([]models.Session, *ape.Error)

	Login(ctx context.Context, request app.LoginRequest) (models.Session, *ape.Error)
	Refresh(ctx context.Context, accountID, sessionID uuid.UUID, request app.RefreshRequest) (models.Session, *ape.Error)
	Logout(ctx context.Context, sessionID uuid.UUID) *ape.Error

	TerminateSessionsByOwner(ctx context.Context, accountID uuid.UUID) *ape.Error
	DeleteSessionByOwner(ctx context.Context, sessionID, initiatorSessionID uuid.UUID) *ape.Error
	TerminateSessionsByAdmin(ctx context.Context, userID uuid.UUID) *ape.Error
	DeleteSessionByAdmin(ctx context.Context, sessionID, initiatorID, initiatorSessionID uuid.UUID) *ape.Error
}

type Controller interface {
	TokenData(w http.ResponseWriter, requestID uuid.UUID, err error)
	ParameterFromURL(w http.ResponseWriter, requestID uuid.UUID, err error, parameter string)
	ResultFromApp(w http.ResponseWriter, requestID uuid.UUID, appErr *ape.Error)
	CheckURLAndJSONResource(w http.ResponseWriter, requestID uuid.UUID, parameter, pointer string)
}

type Handler struct {
	app         App
	controllers Controller
	log         *logrus.Entry
	cfg         config.Config
	google      oauth2.Config
}

func NewHandlers(cfg config.Config, log *logrus.Entry, app *app.App) Handler {
	return Handler{
		app:         app,
		cfg:         cfg,
		google:      config.InitGoogleOAuth(cfg),
		controllers: controllers.NewController(log),
		log:         log,
	}
}
