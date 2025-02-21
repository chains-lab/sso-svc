package domain

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/ape"
	models2 "github.com/recovery-flow/sso-oauth/internal/service/domain/models"
	"github.com/recovery-flow/sso-oauth/internal/service/infra"
	"github.com/recovery-flow/tokens/identity"
	"github.com/sirupsen/logrus"
)

type Domain interface {
	SessionCreate(ctx context.Context, session models2.Session) (*models2.Session, error)
	SessionGet(ctx context.Context, sessionID uuid.UUID) (*models2.Session, error)
	SessionGetForAccount(ctx context.Context, sessionID, accountID uuid.UUID) (*models2.Session, error)
	SessionsListByAccount(ctx context.Context, accountID uuid.UUID) ([]models2.Session, error)

	SessionsTerminate(ctx context.Context, accountID uuid.UUID, excludeSessionID *uuid.UUID) error
	SessionDelete(ctx context.Context, sessionID uuid.UUID) error
	SessionRefresh(ctx context.Context, session models2.Session, role identity.IdnType, IP, client, curToken string) (*string, *string, error)

	Login(ctx context.Context, role identity.IdnType, email, client, IP string) (*string, *string, error)

	AccountCreate(ctx context.Context, acc models2.Account) (*models2.Account, error)
	AccountGet(ctx context.Context, accountID uuid.UUID) (*models2.Account, error)
	AccountGetByEmail(ctx context.Context, email string) (*models2.Account, error)
	AccountUpdateRole(ctx context.Context, accountID uuid.UUID, newRole identity.IdnType) (*models2.Account, error)
}

type domain struct {
	Infra *infra.Infra
	log   *logrus.Logger
}

func NewDomain(cfg *config.Config, log *logrus.Logger) (Domain, error) {
	repo, err := infra.NewDataBase(cfg, log)
	if err != nil {
		return nil, err
	}

	return &domain{
		Infra: repo,
		log:   log,
	}, err
}

func (d *domain) AccountCreate(ctx context.Context, account models2.Account) (*models2.Account, error) {
	res, err := d.Infra.Accounts.Create(ctx, account.Email, account.Role)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *domain) AccountGet(ctx context.Context, accountID uuid.UUID) (*models2.Account, error) {
	account, err := d.Infra.Accounts.GetByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.ErrAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

func (d *domain) AccountGetByEmail(ctx context.Context, email string) (*models2.Account, error) {
	account, err := d.Infra.Accounts.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.ErrAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

func (d *domain) AccountUpdateRole(ctx context.Context, accountID uuid.UUID, newRole identity.IdnType) (*models2.Account, error) {
	account, err := d.Infra.Accounts.UpdateRole(ctx, accountID, newRole)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.ErrAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

func (d *domain) SessionCreate(ctx context.Context, session models2.Session) (*models2.Session, error) {
	ses, err := d.Infra.Sessions.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return ses, nil
}

func (d *domain) SessionGetForAccount(ctx context.Context, sessionID, accountID uuid.UUID) (*models2.Session, error) {
	ses, err := d.Infra.Sessions.GetByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.SessionNotFound
		}
		return nil, err
	}

	if ses.AccountID != accountID {
		return nil, ape.ErrSessionNotBelongToUser
	}

	return ses, nil
}

func (d *domain) SessionGet(ctx context.Context, sessionID uuid.UUID) (*models2.Session, error) {
	ses, err := d.Infra.Sessions.GetByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.SessionNotFound
		}
		return nil, err
	}

	return ses, nil
}

func (d *domain) SessionsListByAccount(ctx context.Context, accountID uuid.UUID) ([]models2.Session, error) {
	ses, err := d.Infra.Sessions.SelectByAccountID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.SessionNotFound
		}
		return nil, err
	}

	return ses, nil
}

func (d *domain) SessionDelete(ctx context.Context, sessionID uuid.UUID) error {
	err := d.Infra.Sessions.Delete(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ape.SessionNotFound
		}
		return err
	}

	return nil
}

func (d *domain) SessionsTerminate(ctx context.Context, accountID uuid.UUID, excludeSessionID *uuid.UUID) error {
	err := d.Infra.Sessions.Terminate(ctx, accountID, excludeSessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ape.SessionNotFound
		}
		return err
	}

	return nil
}

func (d *domain) SessionRefresh(ctx context.Context, session models2.Session, role identity.IdnType, IP, client, curToken string) (*string, *string, error) {
	sessionToken, err := d.Infra.Tokens.DecryptRefresh(session.Token)
	if err != nil {
		return nil, nil, err
	}

	if sessionToken != curToken {
		return nil, nil, ape.ErrTokenInvalid
	}

	refresh, err := d.Infra.Tokens.GenerateRefresh(session.AccountID, session.ID, role)
	if err != nil {
		return nil, nil, err
	}

	access, err := d.Infra.Tokens.GenerateAccess(session.AccountID, session.ID, role)
	if err != nil {
		return nil, nil, err
	}

	refreshCrypto, err := d.Infra.Tokens.EncryptRefresh(refresh)
	if err != nil {
		return nil, nil, err
	}

	_, err = d.Infra.Sessions.UpdateToken(ctx, session.ID, session.AccountID, IP, client, refreshCrypto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ape.SessionNotFound
		}
		return nil, nil, err
	}

	return &access, &refresh, nil
}

func (d *domain) Login(ctx context.Context, role identity.IdnType, email, client, IP string) (*string, *string, error) {
	account, err := d.Infra.Accounts.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			account, err = d.Infra.Accounts.Create(ctx, email, role)
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}

	sessionID := uuid.New()
	refresh, err := d.Infra.Tokens.GenerateRefresh(account.ID, sessionID, role)
	if err != nil {
		return nil, nil, err
	}

	access, err := d.Infra.Tokens.GenerateAccess(account.ID, sessionID, role)
	if err != nil {
		return nil, nil, err
	}

	refreshCrypto, err := d.Infra.Tokens.EncryptRefresh(refresh)
	if err != nil {
		return nil, nil, err
	}

	_, err = d.Infra.Sessions.Create(ctx, models2.Session{
		ID:        sessionID,
		AccountID: account.ID,
		Token:     refreshCrypto,
		IP:        IP,
		Client:    client,
	})
	if err != nil {
		return nil, nil, err
	}

	return &access, &refresh, nil
}
