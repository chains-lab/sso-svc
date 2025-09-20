package handlers

import (
	"context"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/google/uuid"
)

type App interface {
	Register(ctx context.Context, email, password string) error
	RegisterAdmin(
		ctx context.Context,
		initiatorID, initiatorSessionID uuid.UUID,
		params app.RegisterAdminParams,
	) (models.User, error)

	DeleteOwnUser(ctx context.Context, userID, initiatorSessionID uuid.UUID) error
	UpdatePassword(ctx context.Context, userID, SessionID uuid.UUID, currentPassword, newPassword string) error

	ListOwnSessions(
		ctx context.Context,
		userID uuid.UUID,
		pag pagi.Request,
		sort []pagi.SortField,
	) ([]models.Session, pagi.Response, error)
	GetOwnSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error)

	DeleteOwnSessions(ctx context.Context, userID, initiatorSessionID uuid.UUID) error
	DeleteOwnSession(ctx context.Context, userID, initiatorSessionID, sessionID uuid.UUID) error

	RefreshSession(
		ctx context.Context,
		token string,
	) (models.TokensPair, error)

	Login(ctx context.Context, email, password string) (models.TokensPair, error)
	LoginByGoogle(ctx context.Context, email string) (models.TokensPair, error)

	GetUserByID(ctx context.Context, ID uuid.UUID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)

	AdminGetUserSession(ctx context.Context, initiatorID, initiatorSessionID, userID, sessionID uuid.UUID) (models.Session, error)
	AdminListUserSessions(
		ctx context.Context,
		initiatorID, initiatorSessionID, userID uuid.UUID,
		pag pagi.Request,
		sort []pagi.SortField,
	) ([]models.Session, pagi.Response, error)

	AdminDeleteUserSession(ctx context.Context, initiatorID, initiatorSessionID, userID, sessionID uuid.UUID) error
	AdminDeleteUserSessions(ctx context.Context, initiatorID, initiatorSessionID, userID uuid.UUID) error

	AdminGetUser(ctx context.Context, initiatorID, initiatorSessionID, userID uuid.UUID) (models.User, error)
	AdminDeleteUser(
		ctx context.Context,
		initiatorID, initiatorSessionID, userID uuid.UUID,
	) error

	AdminBlockUser(
		ctx context.Context,
		initiatorID, initiatorSessionID, userID uuid.UUID,
	) (models.User, error)
	AdminUnblockUser(
		ctx context.Context,
		initiatorID, initiatorSessionID, userID uuid.UUID,
	) (models.User, error)
}
