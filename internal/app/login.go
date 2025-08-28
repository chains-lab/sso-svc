package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/constant"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) GoogleLogin(ctx context.Context, email, client, ip string) (models.TokensPair, error) {
	user, err := a.users.GetUserByEmail(ctx, email)
	if err != nil {
		return models.TokensPair{}, err
	}

	sessionID := uuid.New()

	access, err := a.jwt.GenerateAccess(user.ID, sessionID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s: %w", user.ID, err))
	}

	refresh, err := a.jwt.GenerateRefresh(user.ID, sessionID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s: %w", user.ID, err))
	}

	refreshCrypto, err := a.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for user %s: %w", user.ID, err))
	}

	ses, err := a.session.CreateUserSession(ctx, user.ID, refreshCrypto, client, ip)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create session for user %s: %w", user.ID, err),
		)
	}

	return models.TokensPair{
		SessionID: ses.ID,
		Refresh:   refresh,
		Access:    access,
	}, nil
}

func (a App) Login(ctx context.Context, email, password, client, ip string) (models.TokensPair, error) {
	user, err := a.GetUserByEmail(ctx, email)
	if err != nil {
		return models.TokensPair{}, err
	}

	if user.Status == constant.UserStatusBlocked {
		return models.TokensPair{}, errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("initator user %s is blocked", user.ID),
		)
	}

	if err = a.users.CheckUserPassword(ctx, user.ID, password); err != nil {
		return models.TokensPair{}, err
	}

	sessionID := uuid.New()

	access, err := a.jwt.GenerateAccess(user.ID, sessionID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s: %w", user.ID, err),
		)
	}

	refresh, err := a.jwt.GenerateRefresh(user.ID, sessionID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s: %w", user.ID, err),
		)
	}

	refreshCrypto, err := a.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for user %s: %w", user.ID, err),
		)
	}

	ses, err := a.session.CreateUserSession(ctx, user.ID, refreshCrypto, client, ip)
	if err != nil {
		return models.TokensPair{}, err
	}

	return models.TokensPair{
		SessionID: ses.ID,
		Refresh:   refresh,
		Access:    access,
	}, nil
}

func (a App) Refresh(ctx context.Context, userID, sessionID uuid.UUID, client, ip, token string) (models.TokensPair, error) {
	user, appErr := a.GetUserByID(ctx, userID)
	if appErr != nil {
		return models.TokensPair{}, appErr
	}

	if user.Status == constant.UserStatusBlocked {
		return models.TokensPair{}, errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("initator user %s is blocked", user.ID),
		)
	}

	session, err := a.session.GetUserSession(ctx, sessionID, user.ID)
	if err != nil {
		return models.TokensPair{}, err
	}

	if session.Client != client {
		return models.TokensPair{}, errx.ErrorSessionClientMismatch.Raise(
			fmt.Errorf("client mismatch"),
		)
	}

	access, err := a.jwt.GenerateAccess(session.UserID, session.ID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s: %w", user.ID, err),
		)
	}

	oldRefresh, err := a.jwt.DecryptRefresh(session.Token)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to decrypt refresh token for user %s: %w", user.ID, err),
		)
	}

	if oldRefresh != token {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("refresh token mismatch"),
		)
	}

	refresh, err := a.session.UpdateToken(ctx, user.ID, session.ID, user.Role, ip)
	if err != nil {
		return models.TokensPair{}, err
	}

	return models.TokensPair{
		SessionID: sessionID,
		Refresh:   refresh,
		Access:    access,
	}, nil
}
