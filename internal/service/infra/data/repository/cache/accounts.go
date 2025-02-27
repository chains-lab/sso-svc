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

type Accounts struct {
	client   *redis.Client
	lifeTime time.Duration
}

func NewAccounts(client *redis.Client, lifetime time.Duration) *Accounts {
	return &Accounts{
		client:   client,
		lifeTime: lifetime,
	}
}

func (a *Accounts) Add(ctx context.Context, account models.Account) error {
	accountKey := fmt.Sprintf("account:id:%s", account.ID)
	emailKey := fmt.Sprintf("account:email:%s", account.Email)

	data := map[string]interface{}{
		"email":      account.Email,
		"role":       string(account.Role),
		"created_at": account.CreatedAt.Format(time.RFC3339),
		"updated_at": account.UpdatedAt.Format(time.RFC3339),
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

func (a *Accounts) GetByID(ctx context.Context, AccountID string) (*models.Account, error) {
	accountKey := fmt.Sprintf("account:id:%s", AccountID)
	vals, err := a.client.HGetAll(ctx, accountKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting account from Redis: %w", err)
	}

	if len(vals) == 0 {
		return nil, fmt.Errorf("account not found, id=%s", AccountID)
	}

	return parseAccount(AccountID, vals)
}

func (a *Accounts) GetByEmail(ctx context.Context, email string) (*models.Account, error) {
	emailKey := fmt.Sprintf("account:email:%s", email)

	accountID, err := a.client.Get(ctx, emailKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting accountID by email: %w", err)
	}

	return a.GetByID(ctx, accountID)
}

func (a *Accounts) Delete(ctx context.Context, AccountID string) error {
	key := fmt.Sprintf("account:id:%s", AccountID)

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

	// Удаляем запись пользователя
	err = a.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting account from Redis: %w", err)
	}

	// Удаляем индекс email → AccountID
	emailKey := fmt.Sprintf("account:email:%s", email)
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

	return account, nil
}
