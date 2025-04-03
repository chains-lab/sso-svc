package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/service/domain/ape"
	"github.com/hs-zavet/sso-oauth/internal/service/domain/models"
	"github.com/hs-zavet/sso-oauth/internal/service/infra"
	"github.com/hs-zavet/tokens/identity"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Domain interface {
	SessionCreate(ctx context.Context, session models.Session) error
	SessionGet(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	SessionGetForAccount(ctx context.Context, sessionID, accountID uuid.UUID) (*models.Session, error)
	SessionsListByAccount(ctx context.Context, accountID uuid.UUID) ([]models.Session, error)

	SessionsTerminate(ctx context.Context, accountID uuid.UUID) error
	SessionDelete(ctx context.Context, sessionID uuid.UUID) error
	SessionRefresh(ctx context.Context, session models.Session, subTypeID *uuid.UUID, role identity.IdnType, IP, client, curToken string) (*string, *string, error)

	SubscriptionUpdate(ctx context.Context, accountID uuid.UUID, subscriptionID *uuid.UUID) error

	Login(ctx context.Context, role identity.IdnType, subTypeID *uuid.UUID, email, client, IP string) (*string, *string, error)

	AccountCreate(ctx context.Context, acc models.Account) (*models.Account, error)
	AccountGet(ctx context.Context, accountID uuid.UUID) (*models.Account, error)
	AccountGetByEmail(ctx context.Context, email string) (*models.Account, error)
	AccountUpdateRole(ctx context.Context, accountID uuid.UUID, newRole identity.IdnType) error
}

type domain struct {
	Infra *infra.Infra
	log   *logrus.Logger
}

func NewDomain(infra *infra.Infra, log *logrus.Logger) (Domain, error) {
	return &domain{
		Infra: infra,
		log:   log,
	}, nil
}

func (d *domain) AccountCreate(ctx context.Context, account models.Account) (*models.Account, error) {
	err := d.Infra.Data.SQL.Accounts.New().Transaction(func(ctx context.Context) error {
		d.log.Debug("Creating account start")

		err := d.Infra.Data.SQL.Accounts.New().Insert(ctx, account)
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return nil, err
	}

	if err = d.Infra.Data.Cache.Accounts.Add(ctx, account); err != nil {
		d.log.WithField("redis", err).Error("failed to add account to cache")
	}

	return &account, nil
}

func (d *domain) AccountGet(ctx context.Context, accountID uuid.UUID) (*models.Account, error) {
	account, err := d.Infra.Data.Cache.Accounts.GetByID(ctx, accountID.String())
	if err != nil && !errors.Is(err, redis.Nil) {
		d.log.WithField("redis", err).Error("failed to get account from cache")
	}
	if account != nil {
		return account, nil
	}

	account, err = d.Infra.Data.SQL.Accounts.New().Filter(map[string]any{"id": accountID}).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.ErrAccountNotFound
		}
		return nil, err
	}

	err = d.Infra.Data.Cache.Accounts.Add(ctx, *account)
	if err != nil {
		d.log.WithField("redis", err).Error("failed to add account to cache")
	}

	return account, nil
}

func (d *domain) AccountGetByEmail(ctx context.Context, email string) (*models.Account, error) {
	account, err := d.Infra.Data.Cache.Accounts.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, redis.Nil) {
		d.log.WithField("redis", err).Error("failed to get account from cache")
	}
	if account != nil {
		return account, nil
	}

	account, err = d.Infra.Data.SQL.Accounts.New().Filter(map[string]any{"email": email}).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.ErrAccountNotFound
		}
		return nil, err
	}

	err = d.Infra.Data.Cache.Accounts.Add(ctx, *account)
	if err != nil {
		d.log.WithField("redis", err).Error("failed to add account to cache")
	}

	return account, nil
}

func (d *domain) AccountUpdateRole(ctx context.Context, accountID uuid.UUID, newRole identity.IdnType) error {
	return d.Infra.Data.SQL.Accounts.New().Transaction(func(ctx context.Context) error {
		err := d.Infra.Data.SQL.Accounts.New().Filter(map[string]any{"id": accountID}).Update(ctx, map[string]any{"role": newRole})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ape.ErrAccountNotFound
			}
			return err
		}

		err = d.Infra.Data.SQL.Sessions.New().Filter(map[string]interface{}{"account_id": accountID}).Delete(ctx)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}

		err = d.Infra.Data.Cache.Accounts.Delete(ctx, accountID.String())
		if err != nil && !errors.Is(err, redis.Nil) {
			if err = d.Infra.Data.Cache.Accounts.Drop(ctx); err != nil {
				d.log.WithField("redis", err).Error("failed to drop account from cache")
			}
			d.log.WithField("redis", err).Error("failed to delete account from cache")
		}

		err = d.Infra.Data.Cache.Sessions.DeleteAllByAccountID(ctx, accountID, nil)
		if err != nil && !errors.Is(err, redis.Nil) {
			if err = d.Infra.Data.Cache.Sessions.Drop(ctx); err != nil {
				d.log.WithField("redis", err).Error("failed to drop sesions from cache")
			}
			d.log.WithField("redis", err).Error("failed to delete sessions from cache")
		}

		return nil
	})
}

func (d *domain) SessionCreate(ctx context.Context, session models.Session) error {
	err := d.Infra.Data.SQL.Sessions.New().Insert(ctx, session)
	if err != nil {
		return err
	}

	err = d.Infra.Data.Cache.Sessions.Add(ctx, session)
	if err != nil {
		d.log.WithField("redis", err).Error("failed to add session to cache")
	}

	return nil
}

func (d *domain) SessionGetForAccount(ctx context.Context, sessionID, accountID uuid.UUID) (*models.Session, error) {
	ses, err := d.Infra.Data.Cache.Sessions.GetByID(ctx, sessionID.String())
	if err != nil || !errors.Is(err, redis.Nil) {
		d.log.WithField("redis", err).Error("failed to get session from cache")
	}
	if ses != nil {
		return ses, nil
	}

	ses, err = d.Infra.Data.SQL.Sessions.New().Filter(map[string]any{"id": sessionID}).Get(ctx)
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

func (d *domain) SessionGet(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	ses, err := d.Infra.Data.Cache.Sessions.GetByID(ctx, sessionID.String())
	if err != nil || !errors.Is(err, redis.Nil) {
		d.log.WithField("redis", err).Error("failed to get session from cache")
	}
	if ses != nil {
		return ses, nil
	}

	ses, err = d.Infra.Data.SQL.Sessions.New().Filter(map[string]any{"id": sessionID.String()}).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.SessionNotFound
		}
		return nil, err
	}

	err = d.Infra.Data.Cache.Sessions.Add(ctx, *ses)
	if err != nil {
		d.log.WithField("redis", err).Error("failed to add session to cache")
	}

	return ses, nil
}

func (d *domain) SessionsListByAccount(ctx context.Context, accountID uuid.UUID) ([]models.Session, error) {
	ses, err := d.Infra.Data.Cache.Sessions.SelectByAccountID(ctx, accountID.String())
	if err != nil || !errors.Is(err, redis.Nil) {
		d.log.WithField("redis", err).Error("failed to get sessions from cache")
	}

	ses, err = d.Infra.Data.SQL.Sessions.New().Filter(map[string]any{"account_id": accountID}).Select(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.SessionNotFound
		}
		return nil, err
	}

	return ses, nil
}

func (d *domain) SessionDelete(ctx context.Context, sessionID uuid.UUID) error {
	err := d.Infra.Data.Cache.Sessions.Delete(ctx, sessionID)
	if err != nil && !errors.Is(err, redis.Nil) {
		if err = d.Infra.Data.Cache.Sessions.Drop(ctx); err != nil {
			d.log.WithField("redis", err).Error("failed to drop session from cache")
		}
		d.log.WithField("redis", err).Error("failed to delete session from cache")
	}

	err = d.Infra.Data.SQL.Sessions.New().Filter(map[string]interface{}{"id": sessionID}).Delete(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ape.SessionNotFound
		}
		return err
	}

	return nil
}

func (d *domain) SessionsTerminate(ctx context.Context, accountID uuid.UUID) error {
	err := d.Infra.Data.Cache.Sessions.DeleteAllByAccountID(ctx, accountID, nil)
	if err != nil && !errors.Is(err, redis.Nil) {
		if err = d.Infra.Data.Cache.Sessions.Drop(ctx); err != nil {
			d.log.WithField("redis", err).Error("failed to drop sessions from cache")
		}
		d.log.WithField("redis", err).Error("failed to delete sessions from cache")
	}

	err = d.Infra.Data.SQL.Sessions.New().Filter(map[string]interface{}{"account_id": accountID}).Delete(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ape.SessionNotFound
		}
		return err
	}

	return nil
}

func (d *domain) SessionRefresh(ctx context.Context, session models.Session, subTypeID *uuid.UUID, role identity.IdnType, IP, client, curToken string) (*string, *string, error) {
	sessionToken, err := d.Infra.Tokens.DecryptRefresh(session.Token)
	if err != nil {
		return nil, nil, err
	}

	if sessionToken != curToken {
		return nil, nil, ape.ErrTokenInvalid
	}

	refresh, err := d.Infra.Tokens.GenerateRefresh(&session.AccountID, &session.ID, subTypeID, role)
	if err != nil {
		return nil, nil, err
	}

	access, err := d.Infra.Tokens.GenerateAccess(&session.AccountID, &session.ID, subTypeID, role)
	if err != nil {
		return nil, nil, err
	}

	refreshCrypto, err := d.Infra.Tokens.EncryptRefresh(refresh)
	if err != nil {
		return nil, nil, err
	}

	err = d.Infra.Data.SQL.Sessions.New().Filter(map[string]interface{}{
		"id":         session.ID,
		"account_id": session.AccountID,
	}).Update(ctx, map[string]interface{}{
		"token":  refreshCrypto,
		"client": client,
		"ip":     IP,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ape.SessionNotFound
		}
		return nil, nil, err
	}

	err = d.Infra.Data.Cache.Sessions.Delete(ctx, session.ID)
	if err != nil && !errors.Is(err, redis.Nil) {
		if err = d.Infra.Data.Cache.Sessions.Drop(ctx); err != nil {
			d.log.WithField("redis", err).Error("failed to drop session from cache")
		}
		d.log.WithField("redis", err).Error("failed to delete session from cache")
	}

	return &access, &refresh, nil
}

func (d *domain) SubscriptionUpdate(ctx context.Context, accountID uuid.UUID, subscriptionID *uuid.UUID) error {
	return d.Infra.Data.SQL.Accounts.New().Transaction(func(ctx context.Context) error {
		err := d.Infra.Data.Cache.Sessions.DeleteAllByAccountID(ctx, accountID, nil)
		if err != nil && !errors.Is(err, redis.Nil) {
			if err = d.Infra.Data.Cache.Sessions.Drop(ctx); err != nil {
				d.log.WithField("redis", err).Error("failed to drop sessions from cache")
			}
			d.log.WithField("redis", err).Error("failed to delete sessions from cache")
		}

		err = d.Infra.Data.Cache.Accounts.Delete(ctx, accountID.String())
		if err != nil && !errors.Is(err, redis.Nil) {
			if err = d.Infra.Data.Cache.Accounts.Drop(ctx); err != nil {
				d.log.WithField("redis", err).Error("failed to drop account from cache")
			}
			d.log.WithField("redis", err).Error("failed to delete account from cache")
		}

		err = d.Infra.Data.SQL.Sessions.New().Filter(map[string]interface{}{"account_id": accountID}).Delete(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ape.SessionNotFound
			}
			return err
		}

		err = d.Infra.Data.SQL.Accounts.New().Filter(map[string]any{
			"id": accountID,
		}).Update(ctx, map[string]any{
			"subscription": subscriptionID,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (d *domain) Login(ctx context.Context, role identity.IdnType, subTypeID *uuid.UUID, email, client, IP string) (*string, *string, error) {
	account, err := d.Infra.Data.SQL.Accounts.New().Filter(map[string]any{
		"email": email,
	}).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			account, err = d.AccountCreate(ctx, models.Account{
				ID:        uuid.New(),
				Email:     email,
				Role:      role,
				UpdatedAt: time.Now(),
				CreatedAt: time.Now(),
			})
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}

	sessionID := uuid.New()
	refresh, err := d.Infra.Tokens.GenerateRefresh(&account.ID, &sessionID, subTypeID, account.Role)
	if err != nil {
		return nil, nil, err
	}

	access, err := d.Infra.Tokens.GenerateAccess(&account.ID, &sessionID, subTypeID, account.Role)
	if err != nil {
		return nil, nil, err
	}

	refreshCrypto, err := d.Infra.Tokens.EncryptRefresh(refresh)
	if err != nil {
		return nil, nil, err
	}

	err = d.Infra.Data.SQL.Sessions.New().Insert(ctx, models.Session{
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
