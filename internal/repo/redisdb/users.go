package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const (
	usersCollection = "users"
	emailIndexKey   = "users:emails"
)

type UserModel struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Role         roles.Role `json:"role"`
	Subscription uuid.UUID  `json:"subscription,omitempty"`
	UpdatedAt    *time.Time `db:"updated_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type Users struct {
	client   *redis.Client
	lifeTime time.Duration
}

// NewUsers создаёт новый инстанс, lifetime задается в минутах.
func NewUsers(client *redis.Client, lifetime int) Users {
	return Users{
		client:   client,
		lifeTime: time.Duration(lifetime) * time.Minute,
	}
}

type CreateUserInput struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Role         roles.Role `json:"role"`
	Subscription uuid.UUID  `json:"subscription,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// Create пытается создать новый аккаунт.
// Записывает данные аккаунта в hash по ключу "users:id:<ID>"
// И в индексный hash "users:emails" устанавливает поле email = userID.
func (a Users) Create(ctx context.Context, input CreateUserInput) error {
	userKey := fmt.Sprintf("%s:id:%s", usersCollection, input.ID.String())

	// Проверяем существование аккаунта по userKey.
	exists, err := a.client.Exists(ctx, userKey).Result()
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if exists > 0 {
		return errors.New("user already exists")
	}

	data := map[string]interface{}{
		"email":        input.Email,
		"role":         string(input.Role),
		"subscription": input.Subscription.String(),
		"created_at":   input.CreatedAt.Format(time.RFC3339),
	}

	if err := a.client.HSet(ctx, userKey, data).Err(); err != nil {
		return fmt.Errorf("error adding user to Redis: %w", err)
	}

	// Устанавливаем индекс: в hash emailIndexKey поле = email, значение = userID.
	if err := a.client.HSet(ctx, emailIndexKey, input.Email, input.ID.String()).Err(); err != nil {
		return fmt.Errorf("error creating email index: %w", err)
	}

	if a.lifeTime > 0 {
		pipe := a.client.Pipeline()
		pipe.Expire(ctx, userKey, a.lifeTime)
		pipe.Expire(ctx, emailIndexKey, a.lifeTime)
		if _, err := pipe.Exec(ctx); err != nil {
			return fmt.Errorf("error setting expiration for keys: %w", err)
		}
	}

	return nil
}

type UserSetInput struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Role         roles.Role `json:"role"`
	Subscription uuid.UUID  `json:"subscription,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// Set перезаписывает данные аккаунта полностью.
// Если аккаунт уже существует, старые данные удаляются, затем создаются новые.
func (a Users) Set(ctx context.Context, input UserSetInput) error {
	userKey := fmt.Sprintf("%s:id:%s", usersCollection, input.ID.String())

	// Удаляем существующую запись по userKey.
	if err := a.client.Del(ctx, userKey).Err(); err != nil {
		//return fmt.Errorf("error deleting existing user key: %w", err)
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

	if err := a.client.HSet(ctx, userKey, data).Err(); err != nil {
		return fmt.Errorf("error setting user in Redis: %w", err)
	}

	if err := a.client.HSet(ctx, emailIndexKey, input.Email, input.ID.String()).Err(); err != nil {
		return fmt.Errorf("error updating email index: %w", err)
	}

	if a.lifeTime > 0 {
		pipe := a.client.Pipeline()
		pipe.Expire(ctx, userKey, a.lifeTime)
		pipe.Expire(ctx, emailIndexKey, a.lifeTime)
		if _, err := pipe.Exec(ctx); err != nil {
			return fmt.Errorf("error setting expiration for keys: %w", err)
		}
	}
	return nil
}

type UserUpdateRequest struct {
	Role         *roles.Role `json:"role"`
	Subscription *uuid.UUID  `json:"subscription,omitempty"`
	UpdatedAt    time.Time   `json:"updated_at,omitempty"`
}

// Update обновляет поля аккаунта. Обновляются данные в основном hash по userKey.
// Индексный hash (emailIndexKey) обновляется только если email изменен – в данном запросе email
// передается отдельно как параметр для идентификации текущего индекса.
func (a Users) Update(ctx context.Context, userID uuid.UUID, input UserUpdateRequest) error {
	userKey := fmt.Sprintf("%s:id:%s", usersCollection, userID.String())

	// Проверяем, существует ли запись по userKey.
	exists, err := a.client.Exists(ctx, userKey).Result()
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("user not found, id=%s", userID.String())
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
	pipe.HSet(ctx, userKey, data)
	// Поскольку индексный hash хранит mapping email -> userID и email здесь не меняется,
	// можно также обновить срок жизни для этого индекса.
	pipe.Expire(ctx, userKey, a.lifeTime)
	pipe.Expire(ctx, emailIndexKey, a.lifeTime)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error updating user in Redis: %w", err)
	}
	return nil
}

// GetByID возвращает аккаунт по userID, используя основной ключ.
func (a Users) GetByID(ctx context.Context, userID string) (UserModel, error) {
	userKey := fmt.Sprintf("%s:id:%s", usersCollection, userID)
	vals, err := a.client.HGetAll(ctx, userKey).Result()
	if err != nil {
		return UserModel{}, fmt.Errorf("error getting user from Redis: %w", err)
	}
	if len(vals) == 0 {
		return UserModel{}, fmt.Errorf("user not found, id=%s", userID)
	}
	return parseUser(userID, vals)
}

// GetByEmail возвращает аккаунт по email, используя индексный hash.
func (a Users) GetByEmail(ctx context.Context, email string) (UserModel, error) {
	userID, err := a.client.HGet(ctx, emailIndexKey, email).Result()
	if err != nil {
		return UserModel{}, fmt.Errorf("error getting userID by email: %w", err)
	}
	return a.GetByID(ctx, userID)
}

// Delete удаляет запись аккаунта и удаляет соответствующую запись из индексного хеша.
func (a Users) Delete(ctx context.Context, userID, email string) error {
	userKey := fmt.Sprintf("%s:id:%s", usersCollection, userID)
	// Не нужно формировать отдельный emailKey, поскольку индекс хранится в одном хеше.

	// Проверяем существование записи.
	exists, err := a.client.Exists(ctx, userKey).Result()
	if err != nil {
		return fmt.Errorf("error checking user existence in Redis: %w", err)
	}
	if exists == 0 {
		return redis.Nil
	}

	pipe := a.client.Pipeline()
	pipe.Del(ctx, userKey)
	// Удаляем поле email из индексного хеша.
	pipe.HDel(ctx, emailIndexKey, email)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error deleting user keys from Redis: %w", err)
	}
	return nil
}

// Drop удаляет все ключи, связанные с аккаунтами, по шаблону.
func (a Users) Drop(ctx context.Context) error {
	pattern := fmt.Sprintf("%s:*", usersCollection)
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

func parseUser(userID string, vals map[string]string) (UserModel, error) {
	createdAt, err := time.Parse(time.RFC3339, vals["created_at"])
	if err != nil {
		return UserModel{}, fmt.Errorf("error parsing created_at: %w", err)
	}

	ID, err := uuid.Parse(userID)
	if err != nil {
		return UserModel{}, fmt.Errorf("error parsing UserID: %w", err)
	}

	subscription, err := uuid.Parse(vals["subscription"])
	if err != nil {
		return UserModel{}, fmt.Errorf("error parsing subscription: %w", err)
	}

	role, err := roles.ParseRole(vals["role"])
	if err != nil {
		return UserModel{}, fmt.Errorf("error parsing role: %w", err)
	}

	res := UserModel{
		ID:           ID,
		Email:        vals["email"],
		Role:         role,
		Subscription: subscription,
		CreatedAt:    createdAt,
	}

	if lastUsed, ok := vals["updated_at"]; ok && lastUsed != "" {
		ua, err := time.Parse(time.RFC3339, lastUsed)
		if err != nil {
			return UserModel{}, fmt.Errorf("error parsing last_used: %w", err)
		}
		res.UpdatedAt = &ua
	}

	return res, nil
}
