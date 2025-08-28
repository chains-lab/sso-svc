package app

import (
	"context"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetOwnSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	return a.session.GetUserSession(ctx, userID, sessionID)
}

func (a App) SelectOwnSessions(
	ctx context.Context,
	userID uuid.UUID,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Session, pagi.Response, error) {
	return a.session.SelectUserSessions(ctx, userID, pag, sort)
}

func (a App) DeleteOwnSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	_, err := a.users.GetInitiatorByID(ctx, userID)
	if err != nil {
		return err
	}

	return a.session.DeleteUserSession(ctx, userID, sessionID)
}

func (a App) DeleteOwnSessions(ctx context.Context, userID uuid.UUID) error {
	_, err := a.users.GetInitiatorByID(ctx, userID)
	if err != nil {
		return err
	}

	return a.session.DeleteUserSessions(ctx, userID)
}

func (a App) DeleteOwn(ctx context.Context, userID uuid.UUID) error {
	_, err := a.users.GetInitiatorByID(ctx, userID)
	if err != nil {
		return err
	}

	return a.users.DeleteUser(ctx, userID)
}
