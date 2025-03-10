package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
	"github.com/recovery-flow/tokens/identity"
	"github.com/redis/go-redis/v9"
)

const accountsCollection = "accounts"

type Accounts interface {
	Add(ctx context.Context, account models.Account) error
	GetByID(ctx context.Context, AccountID string) (*models.Account, error)
	GetByEmail(ctx context.Context, email string) (*models.Account, error)
	Delete(ctx context.Context, AccountID string) error
	Drop(ctx context.Context) error
}

type accounts struct {
	client   *redis.Client
	lifeTime time.Duration
}

func NewAccounts(client *redis.Client, lifetime time.Duration) Accounts {
	return &accounts{
		client:   client,
		lifeTime: lifetime,
	}
}

func (a *accounts) Add(ctx context.Context, account models.Account) error {
	accountKey := fmt.Sprintf("%s:id:%s", accountsCollection, account.ID)
	emailKey := fmt.Sprintf("%s:email:%s", accountsCollection, account.Email)

	data := map[string]interface{}{
		"email":        account.Email,
		"role":         string(account.Role),
		"subscription": account.Subscription,
		"created_at":   account.CreatedAt.Format(time.RFC3339),
		"updated_at":   account.UpdatedAt.Format(time.RFC3339),
	}

	if err := a.client.HSet(ctx, accountKey, data).Err(); err != nil {
		return fmt.Errorf("error adding account to Redis: %w", err)
	}

	if err := a.client.Set(ctx, emailKey, account.ID.String(), 0).Err(); err != nil {
		return fmt.Errorf("error creating email index: %w", err)
	}

	if a.lifeTime > 0 {
		pipe := a.client.Pipeline()
		keys := []string{accountKey, emailKey}
		for _, key := range keys {
			pipe.Expire(ctx, key, a.lifeTime)
		}
		_, err := pipe.Exec(ctx)
		if err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf("error setting expiration for keys: %w", err)
		}
	}

	return nil
}

func (a *accounts) GetByID(ctx context.Context, AccountID string) (*models.Account, error) {
	accountKey := fmt.Sprintf("%s:id:%s", accountsCollection, AccountID)
	vals, err := a.client.HGetAll(ctx, accountKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting account from Redis: %w", err)
	}

	if len(vals) == 0 {
		return nil, fmt.Errorf("account not found, id=%s", AccountID)
	}

	return parseAccount(AccountID, vals)
}

func (a *accounts) GetByEmail(ctx context.Context, email string) (*models.Account, error) {
	emailKey := fmt.Sprintf("%s:email:%s", accountsCollection, email)

	accountID, err := a.client.Get(ctx, emailKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting accountID by email: %w", err)
	}

	return a.GetByID(ctx, accountID)
}

func (a *accounts) Delete(ctx context.Context, AccountID string) error {
	key := fmt.Sprintf("%s:id:%s", accountsCollection, AccountID)

	exists, err := a.client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error checking account existence in Redis: %w", err)
	}

	if exists == 0 {
		return redis.Nil
	}

	email, err := a.client.HGet(ctx, key, "email").Result()
	if err != nil {
		return fmt.Errorf("error getting email for account %s: %w", AccountID, err)
	}

	err = a.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting account from Redis: %w", err)
	}

	emailKey := fmt.Sprintf("%s:email:%s", accountsCollection, email)
	err = a.client.Del(ctx, emailKey).Err()
	if err != nil {
		return fmt.Errorf("error deleting email index: %w", err)
	}

	return nil
}

func parseAccount(AccountID string, vals map[string]string) (*models.Account, error) {
	createdAt, err := time.Parse(time.RFC3339, vals["created_at"])
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, vals["updated_at"])
	if err != nil {
		return nil, fmt.Errorf("error parsing updated_at: %w", err)
	}

	ID, err := uuid.Parse(AccountID)
	if err != nil {
		return nil, fmt.Errorf("error parsing AccountID: %w", err)
	}

	role, err := identity.ParseIdentityType(vals["role"])
	if err != nil {
		return nil, fmt.Errorf("error parsing role: %w", err)
	}

	account := &models.Account{
		ID:        ID,
		Email:     vals["email"],
		Role:      role,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	if exist := vals["subscription"]; exist != "" {
		id, err := uuid.Parse(vals["subscription"])
		if err != nil {
			return nil, fmt.Errorf("error parsing subscription: %w", err)
		}
		account.Subscription = &id
	}

	return account, nil
}

func (a *accounts) Drop(ctx context.Context) error {
	pattern := fmt.Sprintf("%s:*", accountsCollection)
	keys, err := a.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("error fetching keys with pattern %s: %w", pattern, err)
	}
	if len(keys) == 0 {
		return nil
	}
	if err := a.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete keys with pattern %s: %w", pattern, err)
	}
	return nil
}
