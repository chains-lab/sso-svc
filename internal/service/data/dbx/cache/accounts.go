package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
	"github.com/redis/go-redis/v9"
)

type Accounts interface {
	Add(ctx context.Context, account models.Account, ttl time.Duration) error
	GetByID(ctx context.Context, userID string) (*models.Account, error)
	GetByEmail(ctx context.Context, email string) (*models.Account, error)
	Delete(ctx context.Context, userID string) error
}

type accounts struct {
	client   *redis.Client
	LifeTime time.Duration
}

func NewAccounts(client *redis.Client, lifetime time.Duration) Accounts {
	return &accounts{
		client:   client,
		LifeTime: lifetime,
	}
}

func (a *accounts) Add(ctx context.Context, account models.Account, ttl time.Duration) error {
	userKey := fmt.Sprintf("user:id:%s", account.ID)
	emailKey := fmt.Sprintf("user:email:%s", account.Email)

	data := map[string]interface{}{
		"email":      account.Email,
		"role":       account.Role,
		"created_at": account.CreatedAt.Format(time.RFC3339),
		"updated_at": account.UpdatedAt.Format(time.RFC3339),
	}

	err := a.client.HSet(ctx, userKey, data).Err()
	if err != nil {
		return fmt.Errorf("error adding user to Redis: %w", err)
	}

	err = a.client.Set(ctx, emailKey, account.ID.String(), 0).Err()
	if err != nil {
		return fmt.Errorf("error creating email index: %w", err)
	}

	if ttl > 0 {
		_ = a.client.Expire(ctx, userKey, ttl).Err()
		_ = a.client.Expire(ctx, emailKey, ttl).Err()
	}

	return nil
}

func (a *accounts) GetByEmail(ctx context.Context, email string) (*models.Account, error) {
	emailKey := fmt.Sprintf("user:email:%s", email)

	userID, err := a.client.Get(ctx, emailKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting userID by email: %w", err)
	}

	return a.GetByID(ctx, userID)
}

func (a *accounts) GetByID(ctx context.Context, userID string) (*models.Account, error) {
	userKey := fmt.Sprintf("user:id:%s", userID)

	vals, err := a.client.HGetAll(ctx, userKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting user from Redis: %w", err)
	}

	if len(vals) == 0 {
		return nil, fmt.Errorf("user not found, id=%s", userID)
	}

	return parseUser(userID, vals)
}

func (a *accounts) Delete(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:id:%s", userID)

	exists, err := a.client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error checking user existence in Redis: %w", err)
	}

	if exists == 0 {
		return nil
	}

	email, err := a.client.HGet(ctx, key, "email").Result()
	if err != nil {
		return fmt.Errorf("error getting email for user %s: %w", userID, err)
	}

	// Удаляем запись пользователя
	err = a.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting user from Redis: %w", err)
	}

	// Удаляем индекс email → userID
	emailKey := fmt.Sprintf("user:email:%s", email)
	err = a.client.Del(ctx, emailKey).Err()
	if err != nil {
		return fmt.Errorf("error deleting email index: %w", err)
	}

	return nil
}

func parseUser(userID string, vals map[string]string) (*models.Account, error) {
	createdAt, err := time.Parse(time.RFC3339, vals["created_at"])
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, vals["updated_at"])
	if err != nil {
		return nil, fmt.Errorf("error parsing updated_at: %w", err)
	}

	ID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("error parsing userID: %w", err)
	}

	user := &models.Account{
		ID:        ID,
		Email:     vals["email"],
		Role:      vals["role"],
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return user, nil
}
