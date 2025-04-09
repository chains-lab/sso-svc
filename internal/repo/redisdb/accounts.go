package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/tokens/roles"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const (
	accountsCollection = "accounts"
	emailIndexKey      = "accounts:emails"
)

type AccountModel struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Role         roles.Role `json:"role"`
	Subscription uuid.UUID  `json:"subscription,omitempty"`
	UpdatedAt    *time.Time `db:"updated_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type Accounts struct {
	client   *redis.Client
	lifeTime time.Duration
}

// NewAccounts создаёт новый инстанс, lifetime задается в минутах.
func NewAccounts(client *redis.Client, lifetime int) Accounts {
	return Accounts{
		client:   client,
		lifeTime: time.Duration(lifetime) * time.Minute,
	}
}

type CreateAccountInput struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Role         roles.Role `json:"role"`
	Subscription uuid.UUID  `json:"subscription,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// Create пытается создать новый аккаунт.
// Записывает данные аккаунта в hash по ключу "accounts:id:<ID>"
// И в индексный hash "accounts:emails" устанавливает поле email = accountID.
func (a Accounts) Create(ctx context.Context, input CreateAccountInput) error {
	accountKey := fmt.Sprintf("%s:id:%s", accountsCollection, input.ID.String())

	// Проверяем существование аккаунта по accountKey.
	exists, err := a.client.Exists(ctx, accountKey).Result()
	if err != nil {
		return fmt.Errorf("error checking account existence: %w", err)
	}
	if exists > 0 {
		return errors.New("account already exists")
	}

	data := map[string]interface{}{
		"email":        input.Email,
		"role":         string(input.Role),
		"subscription": input.Subscription.String(),
		"created_at":   input.CreatedAt.Format(time.RFC3339),
	}

	if err := a.client.HSet(ctx, accountKey, data).Err(); err != nil {
		return fmt.Errorf("error adding account to Redis: %w", err)
	}

	// Устанавливаем индекс: в hash emailIndexKey поле = email, значение = accountID.
	if err := a.client.HSet(ctx, emailIndexKey, input.Email, input.ID.String()).Err(); err != nil {
		return fmt.Errorf("error creating email index: %w", err)
	}

	if a.lifeTime > 0 {
		pipe := a.client.Pipeline()
		pipe.Expire(ctx, accountKey, a.lifeTime)
		pipe.Expire(ctx, emailIndexKey, a.lifeTime)
		if _, err := pipe.Exec(ctx); err != nil {
			return fmt.Errorf("error setting expiration for keys: %w", err)
		}
	}

	return nil
}

type AccountSetInput struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Role         roles.Role `json:"role"`
	Subscription uuid.UUID  `json:"subscription,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// Set перезаписывает данные аккаунта полностью.
// Если аккаунт уже существует, старые данные удаляются, затем создаются новые.
func (a Accounts) Set(ctx context.Context, input AccountSetInput) error {
	accountKey := fmt.Sprintf("%s:id:%s", accountsCollection, input.ID.String())

	// Удаляем существующую запись по accountKey.
	if err := a.client.Del(ctx, accountKey).Err(); err != nil {
		//return fmt.Errorf("error deleting existing account key: %w", err)
	}

	// Обновляем индекс email: в emailIndexKey ставим для поля input.Email значение input.ID.
	if err := a.client.HDel(ctx, emailIndexKey, input.Email).Err(); err != nil {
		// Если ключа нет, можно проигнорировать ошибку.
	}

	data := map[string]interface{}{
		"email":        input.Email,
		"role":         string(input.Role),
		"subscription": input.Subscription.String(),
		"created_at":   input.CreatedAt.Format(time.RFC3339),
	}
	if input.UpdatedAt != nil {
		data["updated_at"] = input.UpdatedAt.Format(time.RFC3339)
	}

	if err := a.client.HSet(ctx, accountKey, data).Err(); err != nil {
		return fmt.Errorf("error setting account in Redis: %w", err)
	}

	if err := a.client.HSet(ctx, emailIndexKey, input.Email, input.ID.String()).Err(); err != nil {
		return fmt.Errorf("error updating email index: %w", err)
	}

	if a.lifeTime > 0 {
		pipe := a.client.Pipeline()
		pipe.Expire(ctx, accountKey, a.lifeTime)
		pipe.Expire(ctx, emailIndexKey, a.lifeTime)
		if _, err := pipe.Exec(ctx); err != nil {
			return fmt.Errorf("error setting expiration for keys: %w", err)
		}
	}
	return nil
}

type AccountUpdateRequest struct {
	Role         *roles.Role `json:"role"`
	Subscription *uuid.UUID  `json:"subscription,omitempty"`
	UpdatedAt    time.Time   `json:"updated_at,omitempty"`
}

// Update обновляет поля аккаунта. Обновляются данные в основном hash по accountKey.
// Индексный hash (emailIndexKey) обновляется только если email изменен – в данном запросе email
// передается отдельно как параметр для идентификации текущего индекса.
func (a Accounts) Update(ctx context.Context, accountID uuid.UUID, input AccountUpdateRequest) error {
	accountKey := fmt.Sprintf("%s:id:%s", accountsCollection, accountID.String())

	// Проверяем, существует ли запись по accountKey.
	exists, err := a.client.Exists(ctx, accountKey).Result()
	if err != nil {
		return fmt.Errorf("error checking account existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("account not found, id=%s", accountID.String())
	}

	data := make(map[string]interface{})
	if input.Role != nil {
		data["role"] = string(*input.Role)
	}
	if input.Subscription != nil {
		data["subscription"] = input.Subscription.String()
	}
	// Обновляем временную метку.
	data["updated_at"] = input.UpdatedAt.Format(time.RFC3339)

	pipe := a.client.Pipeline()
	pipe.HSet(ctx, accountKey, data)
	// Поскольку индексный hash хранит mapping email -> accountID и email здесь не меняется,
	// можно также обновить срок жизни для этого индекса.
	pipe.Expire(ctx, accountKey, a.lifeTime)
	pipe.Expire(ctx, emailIndexKey, a.lifeTime)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error updating account in Redis: %w", err)
	}
	return nil
}

// GetByID возвращает аккаунт по accountID, используя основной ключ.
func (a Accounts) GetByID(ctx context.Context, accountID string) (AccountModel, error) {
	accountKey := fmt.Sprintf("%s:id:%s", accountsCollection, accountID)
	vals, err := a.client.HGetAll(ctx, accountKey).Result()
	if err != nil {
		return AccountModel{}, fmt.Errorf("error getting account from Redis: %w", err)
	}
	if len(vals) == 0 {
		return AccountModel{}, fmt.Errorf("account not found, id=%s", accountID)
	}
	return parseAccount(accountID, vals)
}

// GetByEmail возвращает аккаунт по email, используя индексный hash.
func (a Accounts) GetByEmail(ctx context.Context, email string) (AccountModel, error) {
	accountID, err := a.client.HGet(ctx, emailIndexKey, email).Result()
	if err != nil {
		return AccountModel{}, fmt.Errorf("error getting accountID by email: %w", err)
	}
	return a.GetByID(ctx, accountID)
}

// Delete удаляет запись аккаунта и удаляет соответствующую запись из индексного хеша.
func (a Accounts) Delete(ctx context.Context, accountID, email string) error {
	accountKey := fmt.Sprintf("%s:id:%s", accountsCollection, accountID)
	// Не нужно формировать отдельный emailKey, поскольку индекс хранится в одном хеше.

	// Проверяем существование записи.
	exists, err := a.client.Exists(ctx, accountKey).Result()
	if err != nil {
		return fmt.Errorf("error checking account existence in Redis: %w", err)
	}
	if exists == 0 {
		return redis.Nil
	}

	pipe := a.client.Pipeline()
	pipe.Del(ctx, accountKey)
	// Удаляем поле email из индексного хеша.
	pipe.HDel(ctx, emailIndexKey, email)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error deleting account keys from Redis: %w", err)
	}
	return nil
}

// Drop удаляет все ключи, связанные с аккаунтами, по шаблону.
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

func parseAccount(accountID string, vals map[string]string) (AccountModel, error) {
	createdAt, err := time.Parse(time.RFC3339, vals["created_at"])
	if err != nil {
		return AccountModel{}, fmt.Errorf("error parsing created_at: %w", err)
	}

	ID, err := uuid.Parse(accountID)
	if err != nil {
		return AccountModel{}, fmt.Errorf("error parsing AccountID: %w", err)
	}

	subscription, err := uuid.Parse(vals["subscription"])
	if err != nil {
		return AccountModel{}, fmt.Errorf("error parsing subscription: %w", err)
	}

	role, err := roles.ParseRole(vals["role"])
	if err != nil {
		return AccountModel{}, fmt.Errorf("error parsing role: %w", err)
	}

	res := AccountModel{
		ID:           ID,
		Email:        vals["email"],
		Role:         role,
		Subscription: subscription,
		CreatedAt:    createdAt,
	}

	if lastUsed, ok := vals["updated_at"]; ok && lastUsed != "" {
		ua, err := time.Parse(time.RFC3339, lastUsed)
		if err != nil {
			return AccountModel{}, fmt.Errorf("error parsing last_used: %w", err)
		}
		res.UpdatedAt = &ua
	}

	return res, nil
}
