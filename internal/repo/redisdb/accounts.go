package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const accountsCollection = "accounts"

type AccountModel struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Role         string     `json:"role"`
	Subscription uuid.UUID  `json:"subscription,omitempty"`
	UpdatedAt    *time.Time `db:"updated_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type Accounts struct {
	client   *redis.Client
	lifeTime time.Duration
}

func NewAccounts(client *redis.Client, lifetime int) Accounts {
	return Accounts{
		client:   client,
		lifeTime: time.Duration(lifetime) * time.Minute,
	}
}

type InsertAccountInput struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Role         string     `json:"role"`
	Subscription uuid.UUID  `json:"subscription,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (a Accounts) Create(ctx context.Context, input InsertAccountInput) error {
	accountKey := fmt.Sprintf("%s:id:%s", accountsCollection, input.ID)
	emailKey := fmt.Sprintf("%s:email:%s", accountsCollection, input.Email)

	exists, err := a.client.Exists(ctx, accountKey).Result()
	if err != nil {
		return fmt.Errorf("error checking account existence: %w", err)
	}
	if exists > 0 {
		return errors.New("account already exists")
	}

	data := map[string]interface{}{
		"email":        input.Email,
		"role":         input.Role,
		"subscription": input.Subscription.String(),
		"created_at":   input.CreatedAt.Format(time.RFC3339),
	}
	if input.UpdatedAt != nil {
		data["updated_at"] = input.UpdatedAt.Format(time.RFC3339)
	}

	if err := a.client.HSet(ctx, accountKey, data).Err(); err != nil {
		return fmt.Errorf("error adding account to Redis: %w", err)
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
		if _, err := pipe.Exec(ctx); err != nil {
			return fmt.Errorf("error setting expiration for keys: %w", err)
		}
	}

	return nil
}

func (a Accounts) Set(ctx context.Context, input InsertAccountInput) error {
	accountKey := fmt.Sprintf("%s:id:%s", accountsCollection, input.ID)
	emailKey := fmt.Sprintf("%s:email:%s", accountsCollection, input.Email)

	if err := a.client.Del(ctx, accountKey, emailKey).Err(); err != nil {
		return fmt.Errorf("error deleting existing keys: %w", err)
	}

	data := map[string]interface{}{
		"email":        input.Email,
		"role":         input.Role,
		"subscription": input.Subscription.String(),
		"created_at":   input.CreatedAt.Format(time.RFC3339),
	}
	if input.UpdatedAt != nil {
		data["updated_at"] = input.UpdatedAt.Format(time.RFC3339)
	}

	if err := a.client.HSet(ctx, accountKey, data).Err(); err != nil {
		return fmt.Errorf("error setting account in Redis: %w", err)
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
		if _, err := pipe.Exec(ctx); err != nil {
			return fmt.Errorf("error setting expiration for keys: %w", err)
		}
	}

	return nil
}

type AccountUpdateRequest struct {
	Role         *string    `json:"role"`
	Subscription *uuid.UUID `json:"subscription,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at,omitempty"`
}

func (a Accounts) Update(ctx context.Context, accountID uuid.UUID, input AccountUpdateRequest) error {
	accountKey := fmt.Sprintf("%s:id:%s", accountsCollection, accountID)

	exists, err := a.client.Exists(ctx, accountKey).Result()
	if err != nil {
		return fmt.Errorf("error checking account existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("account not found, id=%s", accountID)
	}

	data := make(map[string]interface{})

	if input.Role != nil {
		data["role"] = input.Role
	}

	if input.Subscription != nil {
		data["subscription"] = input.Subscription
	}

	data["updated_at"] = input.UpdatedAt.Format(time.RFC3339)

	if err := a.client.HSet(ctx, accountKey, data).Err(); err != nil {
		return fmt.Errorf("error updating account in Redis: %w", err)
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

	ID, err := uuid.Parse(AccountID)
	if err != nil {
		return AccountModel{}, fmt.Errorf("error parsing AccountID: %w", err)
	}

	//role, err := roles.ParseRole(vals["role"])
	//if err != nil {
	//	return AccountModel{}, fmt.Errorf("error parsing role: %w", err)
	//}

	subscription, err := uuid.Parse(vals["subscription"])
	if err != nil {
		return AccountModel{}, fmt.Errorf("error parsing subscription: %w", err)
	}

	res := AccountModel{
		ID:           ID,
		Email:        vals["email"],
		Role:         vals["role"],
		Subscription: subscription,
		CreatedAt:    createdAt,
	}

	if lastUsed, ok := vals["updated_at"]; ok {
		ua, err := time.Parse(time.RFC3339, lastUsed)
		if err != nil {
			return AccountModel{}, fmt.Errorf("error parsing last_used: %w", err)
		}
		res.UpdatedAt = &ua
	}

	return res, nil
}
