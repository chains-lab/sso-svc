package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/rerabbit"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/ape"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
	"github.com/recovery-flow/sso-oauth/internal/service/infra"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/events/rabbit/amqpconfig"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/events/rabbit/evebody"
	"github.com/recovery-flow/tokens/identity"
	"github.com/sirupsen/logrus"
)

type Domain interface {
	SessionCreate(ctx context.Context, session models.Session) (*models.Session, error)
	SessionGet(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	SessionGetForAccount(ctx context.Context, sessionID, accountID uuid.UUID) (*models.Session, error)
	SessionsListByAccount(ctx context.Context, accountID uuid.UUID) ([]models.Session, error)

	SessionsTerminate(ctx context.Context, accountID uuid.UUID, excludeSessionID *uuid.UUID) error
	SessionDelete(ctx context.Context, sessionID uuid.UUID) error
	SessionRefresh(ctx context.Context, session models.Session, role identity.IdnType, IP, client, curToken string) (*string, *string, error)

	Login(ctx context.Context, role identity.IdnType, email, client, IP string) (*string, *string, error)

	AccountCreate(ctx context.Context, acc models.Account) (*models.Account, error)
	AccountGet(ctx context.Context, accountID uuid.UUID) (*models.Account, error)
	AccountGetByEmail(ctx context.Context, email string) (*models.Account, error)
	AccountUpdateRole(ctx context.Context, accountID uuid.UUID, newRole identity.IdnType) (*models.Account, error)
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
	res, err := d.Infra.Data.Accounts.Create(ctx, account.Email, account.Role)
	if err != nil {
		return nil, err
	}

	eventBody := evebody.AccountCreated{
		Event:     amqpconfig.AccountCreateKey,
		AccountID: res.ID.String(),
		Email:     res.Email,
		Role:      string(res.Role),
		Timestamp: time.Now().UTC(),
	}

	bodyBytes, err := json.Marshal(eventBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	publishOpts := rerabbit.PublishOptions{
		Exchange:      amqpconfig.AccountSSOExchange,
		RoutingKey:    amqpconfig.AccountCreateKey,
		Mandatory:     false,
		Immediate:     false,
		ContentType:   "application/json",
		DeliveryMode:  2, // 2 = Persistent
		Headers:       nil,
		CorrelationID: "",
		ReplyTo:       "",
		Body:          bodyBytes,
	}

	err = d.Infra.Rabbit.Publish(ctx, publishOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to publish event: %w", err)
	}
	return res, nil
}

func (d *domain) AccountGet(ctx context.Context, accountID uuid.UUID) (*models.Account, error) {
	account, err := d.Infra.Data.Accounts.GetByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.ErrAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

func (d *domain) AccountGetByEmail(ctx context.Context, email string) (*models.Account, error) {
	account, err := d.Infra.Data.Accounts.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.ErrAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

func (d *domain) AccountUpdateRole(ctx context.Context, accountID uuid.UUID, newRole identity.IdnType) (*models.Account, error) {
	account, err := d.Infra.Data.Accounts.UpdateRole(ctx, accountID, newRole)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.ErrAccountNotFound
		}
		return nil, err
	}

	// Формируем событие обновления роли
	eventBody := evebody.AccountRoleUpdated{
		Event:     amqpconfig.AccountUpdateRoleKey,
		AccountID: account.ID.String(),
		Role:      string(newRole),
		Timestamp: time.Now().UTC(),
	}

	// Маршаллинг события в JSON
	bodyBytes, err := json.Marshal(eventBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	publishOpts := rerabbit.PublishOptions{
		Exchange:      amqpconfig.AccountSSOExchange,
		RoutingKey:    amqpconfig.AccountUpdateRoleKey,
		Mandatory:     false,
		Immediate:     false,
		ContentType:   "application/json",
		DeliveryMode:  2, // Persistent
		Headers:       nil,
		CorrelationID: "",
		ReplyTo:       "",
		Body:          bodyBytes,
	}

	err = d.Infra.Rabbit.Publish(ctx, publishOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to publish event: %w", err)
	}

	return account, nil
}

func (d *domain) SessionCreate(ctx context.Context, session models.Session) (*models.Session, error) {
	ses, err := d.Infra.Data.Sessions.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return ses, nil
}

func (d *domain) SessionGetForAccount(ctx context.Context, sessionID, accountID uuid.UUID) (*models.Session, error) {
	ses, err := d.Infra.Data.Sessions.GetByID(ctx, sessionID)
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
	ses, err := d.Infra.Data.Sessions.GetByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.SessionNotFound
		}
		return nil, err
	}

	return ses, nil
}

func (d *domain) SessionsListByAccount(ctx context.Context, accountID uuid.UUID) ([]models.Session, error) {
	ses, err := d.Infra.Data.Sessions.SelectByAccountID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ape.SessionNotFound
		}
		return nil, err
	}

	return ses, nil
}

func (d *domain) SessionDelete(ctx context.Context, sessionID uuid.UUID) error {
	err := d.Infra.Data.Sessions.Delete(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ape.SessionNotFound
		}
		return err
	}

	return nil
}

func (d *domain) SessionsTerminate(ctx context.Context, accountID uuid.UUID, excludeSessionID *uuid.UUID) error {
	err := d.Infra.Data.Sessions.Terminate(ctx, accountID, excludeSessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ape.SessionNotFound
		}
		return err
	}

	return nil
}

func (d *domain) SessionRefresh(ctx context.Context, session models.Session, role identity.IdnType, IP, client, curToken string) (*string, *string, error) {
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

	_, err = d.Infra.Data.Sessions.UpdateToken(ctx, session.ID, session.AccountID, IP, client, refreshCrypto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ape.SessionNotFound
		}
		return nil, nil, err
	}

	return &access, &refresh, nil
}

func (d *domain) Login(ctx context.Context, role identity.IdnType, email, client, IP string) (*string, *string, error) {
	account, err := d.Infra.Data.Accounts.GetByEmail(ctx, email)
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
	refresh, err := d.Infra.Tokens.GenerateRefresh(account.ID, sessionID, account.Role)
	if err != nil {
		return nil, nil, err
	}

	access, err := d.Infra.Tokens.GenerateAccess(account.ID, sessionID, account.Role)
	if err != nil {
		return nil, nil, err
	}

	refreshCrypto, err := d.Infra.Tokens.EncryptRefresh(refresh)
	if err != nil {
		return nil, nil, err
	}

	_, err = d.Infra.Data.Sessions.Create(ctx, models.Session{
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
