package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) AdminGetUserSession(ctx context.Context, initiatorID, initiatorSessionID, userID, sessionID uuid.UUID) (models.Session, error) {
	_, _, err := a.users.CompareRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return models.Session{}, err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return models.Session{}, err
	}

	session, appErr := a.session.GetForUser(ctx, userID, sessionID)
	if appErr != nil {
		return models.Session{}, appErr
	}

	return session, nil
}

func (a App) AdminListUserSessions(
	ctx context.Context,
	initiatorID, initiatorSessionID, userID uuid.UUID,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Session, pagi.Response, error) {
	_, _, err := a.users.CompareRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return nil, pagi.Response{}, err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return nil, pagi.Response{}, err
	}

	return a.session.ListForUser(ctx, userID, pag, sort)
}

func (a App) AdminDeleteUserSession(ctx context.Context, initiatorID, initiatorSessionID, userID, sessionID uuid.UUID) error {
	_, _, err := a.users.CompareRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return err
	}

	session, appErr := a.session.GetForUser(ctx, userID, sessionID)
	if appErr != nil {
		return appErr
	}

	if session.UserID != userID {
		return errx.ErrorSessionNotFound.Raise(
			fmt.Errorf("session %s does not belong to user %s", sessionID, userID),
		)
	}

	appErr = a.session.DeleteOneForUser(ctx, userID, session.ID)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (a App) AdminDeleteUserSessions(ctx context.Context, initiatorID, initiatorSessionID, userID uuid.UUID) error {
	_, _, err := a.users.CompareRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return err
	}

	err = a.transaction(func(ctx context.Context) error {
		err = a.session.DeleteAllForUser(ctx, userID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
