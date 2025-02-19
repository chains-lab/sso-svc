package domain

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/models"
	"github.com/recovery-flow/sso-oauth/internal/service/infra"
	"github.com/sirupsen/logrus"
)

type Domain interface {
	SessionGet(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	SessionGetForUser(ctx context.Context, sessionID, userID uuid.UUID) (*models.Session, error)
	SessionCreate(ctx context.Context, session models.Session) (*models.Session, error)
	SessionsListByUser(ctx context.Context, userID uuid.UUID) ([]models.Session, error)

	SessionsTerminate(ctx context.Context, userID uuid.UUID, excludeSessionID *uuid.UUID) error
	SessionDelete(ctx context.Context, sessionID uuid.UUID) error
	SessionRefresh(
		ctx context.Context,
		session *models.Session,
		curToken string,
		role roles.UserRole,
		IP string,
	) (string, string, error)
	SessionLogin(
		ctx context.Context,
		email string,
		role roles.UserRole,
		IP string,
	) (string, string, error)

	AccountCreate(ctx context.Context, acc models.Account) (*models.Account, error)
	AccountGet(ctx context.Context, accountID uuid.UUID) (*models.Account, error)
	AccountGetByEmail(ctx context.Context, email string) (*models.Account, error)
	AccountUpdateRole(ctx context.Context, accountID uuid.UUID, newRole string) (*models.Account, error)
}

type domain struct {
	Infra *infra.Infra
	log   *logrus.Logger
}

func NewDomain(cfg *config.Config, logger *logrus.Logger) (Domain, error) {
	repo, err := infra.NewDataBase(cfg)
	if err != nil {
		return nil, err
	}

	return &domain{
		Infra: repo,
		log:   logger,
	}, err
}

func (d *domain) AccountGet(ctx context.Context, accountID uuid.UUID) (*models.Account, error) {
	user, err := d.Infra.Accounts.GetByID(ctx, accountID)
	if err != nil {
		d.log.Errorf("Failed to retrieve user: %v", err)
		return nil, err
	}

	return user, nil
}

func (d *domain) AccountGetByEmail(ctx context.Context, email string) (*models.Account, error) {
	user, err := d.Infra.Accounts.GetByEmail(ctx, email)
	if err != nil {
		d.log.Errorf("Failed to retrieve user: %v", err)
		return nil, err
	}

	return user, nil
}

func (d *domain) AccountCreate(ctx context.Context, account models.Account) (*models.Account, error) {
	user, err := d.Infra.Accounts.Create(ctx, account.Email, account.Role)
	if err != nil {
		d.log.Errorf("Failed to create user: %v", err)
		return nil, err
	}

	return user, nil
}

func (d *domain) AccountUpdateRole(ctx context.Context, accountID uuid.UUID, newRole string) (*models.Account, error) {
	accRole, err := roles.ParseUserRole(newRole)
	if err != nil {
		return nil, err
	}
	user, err := d.Infra.Accounts.UpdateRole(ctx, accountID, accRole)
	if err != nil {
		d.log.Errorf("Failed to update user role: %v", err)
		return nil, err
	}

	return user, nil
}

func (d *domain) SessionCreate(ctx context.Context, session models.Session) (*models.Session, error) {
	ses, err := d.Infra.Sessions.Create(ctx, session)
	if err != nil {
		d.log.Errorf("Failed to create sessions: %v", err)
		return nil, fmt.Errorf("failed to create sessions: %w", err)
	}

	return ses, nil
}

func (d *domain) SessionGetForUser(ctx context.Context, sessionID, userID uuid.UUID) (*models.Session, error) {
	ses, err := d.Infra.Sessions.GetByID(ctx, sessionID)
	if err != nil {
		d.log.Errorf("Failed to retrieve user sessions: %v", err)
		return nil, fmt.Errorf("failed to retrieve user sessions: %w", err)
	}

	if ses.UserID != userID {
		d.log.Debugf("Sessions doesn't belong to user")
		return nil, problems.Forbidden("Sessions doesn't belong to user")
	}

	return ses, nil
}

func (d *domain) SessionGet(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	ses, err := d.Infra.Sessions.GetByID(ctx, sessionID)
	if err != nil {
		d.log.Errorf("Failed to retrieve user sessions: %v", err)
		return nil, fmt.Errorf("failed to retrieve user sessions: %w", err)
	}

	return ses, nil
}

func (d *domain) SessionsListByUser(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	ses, err := d.Infra.Sessions.SelectByUserID(ctx, userID)
	if err != nil {
		d.log.Errorf("Failed to retrieve user sessions: %v", err)
		return nil, fmt.Errorf("failed to retrieve user sessions: %w", err)
	}

	return ses, nil
}

func (d *domain) SessionDelete(ctx context.Context, sessionID uuid.UUID) error {
	err := d.Infra.Sessions.Delete(ctx, sessionID)
	if err != nil {
		d.log.Errorf("Failed to delete sessions: %v", err)
		return fmt.Errorf("failed to delete sessions: %w", err)
	}

	return nil
}

func (d *domain) SessionsTerminate(ctx context.Context, userID uuid.UUID, excludeSessionID *uuid.UUID) error {
	err := d.Infra.Sessions.Terminate(ctx, userID, excludeSessionID)
	if err != nil {
		d.log.Errorf("Failed to terminate user sessions: %v", err)
		return fmt.Errorf("failed to terminate user sessions: %w", err)
	}

	return nil
}

func (d *domain) SessionRefresh(
	ctx context.Context,
	session *models.Session,
	curToken string,
	role roles.UserRole,
	IP string,
) (string, string, error) {
	sessionToken, err := d.Infra.Tokens.DecryptRefresh(session.Token)
	if err != nil {
		d.log.Errorf("Failed to decrypt refresh token: %v", err)
		return "", "", problems.InternalError()
	}

	if sessionToken != curToken {
		d.log.Debugf("Invalid refresh token")
		return "", "", problems.Unauthorized("Invalid refresh token")
	}

	refresh, err := d.Infra.Tokens.GenerateRefresh(session.UserID, session.ID, role)
	if err != nil {
		d.log.Errorf("Failed to generate refresh token: %v", err)
		return "", "", problems.InternalError()
	}

	access, err := d.Infra.Tokens.GenerateAccess(session.UserID, session.ID, role)
	if err != nil {
		d.log.Errorf("Failed to generate access token: %v", err)
		return "", "", problems.InternalError()
	}

	_, err = d.Infra.Sessions.UpdateToken(ctx, session.ID, refresh, IP)
	if err != nil {
		d.log.Errorf("Failed to update session token: %v", err)
		return "", "", problems.InternalError()
	}

	return access, refresh, nil
}

func (d *domain) SessionLogin(
	ctx context.Context,
	email string,
	role roles.UserRole,
	IP string,
) (string, string, error) {
	var accountID, sessionID uuid.UUID

	account, err := d.AccountGetByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			accountID = uuid.New()

			account, err = d.AccountCreate(ctx, models.Account{
				ID:        accountID,
				Email:     email,
				Role:      role,
				UpdatedAt: time.Now(),
				CreatedAt: time.Now(),
			})
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					d.log.Errorf("error creating user: %v", err)
					return "", "", problems.InternalError()
				}
			}
		}
	}

	sessionID = uuid.New()
	refresh, err := d.Infra.Tokens.GenerateRefresh(account.ID, sessionID, role)
	if err != nil {
		d.log.Errorf("Failed to generate refresh token: %v", err)
		return "", "", problems.InternalError()
	}

	access, err := d.Infra.Tokens.GenerateAccess(account.ID, sessionID, role)
	if err != nil {
		d.log.Errorf("Failed to generate access token: %v", err)
		return "", "", problems.InternalError()
	}

	_, err = d.Infra.Sessions.UpdateToken(ctx, sessionID, refresh, IP)
	if err != nil {
		d.log.Errorf("Failed to update session token: %v", err)
		return "", "", problems.InternalError()
	}

	return access, refresh, nil
}
