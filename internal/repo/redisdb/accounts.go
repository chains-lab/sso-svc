package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/tokens/identity"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const accountsCollection = "accounts"

type AccountModel struct {
	ID           uuid.UUID        `json:"id"`
	Email        string           `json:"email"`
	Role         identity.IdnType `json:"role"`
	Subscription uuid.UUID        `json:"subscription,omitempty"`
	UpdatedAt    time.Time        `json:"updated_at"`
	CreatedAt    time.Time        `json:"created_at"`
}

type Accounts struct {
	client   *redis.Client
	lifeTime time.Duration
}

func NewAccounts(client *redis.Client, lifetime time.Duration) Accounts {
	return Accounts{
		client:   client,
		lifeTime: lifetime,
	}
}

type AddAccountInput struct {
	ID           uuid.UUID        `json:"id"`
	Email        string           `json:"email"`
	Role         identity.IdnType `json:"role"`
	Subscription uuid.UUID        `json:"subscription,omitempty"`
	UpdatedAt    time.Time        `json:"updated_at"`
	CreatedAt    time.Time        `json:"created_at"`
}

func (a Accounts) Add(ctx context.Context, input AddAccountInput) error {
	accountKey := fmt.Sprintf("%s:id:%s", accountsCollection, input.ID)
	emailKey := fmt.Sprintf("%s:email:%s", accountsCollection, input.Email)

	data := map[string]interface{}{
		"email":        input.Email,
		"role":         string(input.Role),
		"subscription": input.Subscription.String(),
		"created_at":   input.CreatedAt.Format(time.RFC3339),
		"updated_at":   input.UpdatedAt.Format(time.RFC3339),
	}

	if err := a.client.HSet(ctx, accountKey, data).Err(); err != nil {
		return fmt.Errorf("error adding input to Redis: %w", err)
	}

	if err := a.client.Set(ctx, emailKey, input.ID.String(), 0).Err(); err != nil {
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

func (a Accounts) GetByID(ctx context.Context, AccountID string) (AccountModel, error) {
	accountKey := fmt.Sprintf("%s:id:%s", accountsCollection, AccountID)
	vals, err := a.client.HGetAll(ctx, accountKey).Result()
	if err != nil {
		return AccountModel{}, fmt.Errorf("error getting account from Redis: %w", err)
	}

	if len(vals) == 0 {
		return AccountModel{}, fmt.Errorf("account not found, id=%s", AccountID)
	}

	return parseAccount(AccountID, vals)
}

func (a Accounts) GetByEmail(ctx context.Context, email string) (AccountModel, error) {
	emailKey := fmt.Sprintf("%s:email:%s", accountsCollection, email)

	accountID, err := a.client.Get(ctx, emailKey).Result()
	if err != nil {
		return AccountModel{}, fmt.Errorf("error getting accountID by email: %w", err)
	}

	return a.GetByID(ctx, accountID)
}

func (a Accounts) Delete(ctx context.Context, AccountID string) error {
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

func (a Accounts) Drop(ctx context.Context) error {
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

func parseAccount(AccountID string, vals map[string]string) (AccountModel, error) {
	createdAt, err := time.Parse(time.RFC3339, vals["created_at"])
	if err != nil {
		return AccountModel{}, fmt.Errorf("error parsing created_at: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, vals["updated_at"])
	if err != nil {
		return AccountModel{}, fmt.Errorf("error parsing updated_at: %w", err)
	}

	ID, err := uuid.Parse(AccountID)
	if err != nil {
		return AccountModel{}, fmt.Errorf("error parsing AccountID: %w", err)
	}

	role, err := identity.ParseIdentityType(vals["role"])
	if err != nil {
		return AccountModel{}, fmt.Errorf("error parsing role: %w", err)
	}

	subscription, err := uuid.Parse(vals["subscription"])
	if err != nil {
		return AccountModel{}, fmt.Errorf("error parsing subscription: %w", err)
	}

	return AccountModel{
		ID:           ID,
		Email:        vals["email"],
		Role:         role,
		Subscription: subscription,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}
