package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const sessionsCollection = "sessions"

type SessionModel struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	Client    string    `json:"client"`
	LastUsed  time.Time `json:"last_used,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Sessions struct {
	client   *redis.Client
	lifeTime time.Duration
}

func NewSessions(client *redis.Client, lifetime int) Sessions {
	return Sessions{
		client:   client,
		lifeTime: time.Duration(lifetime) * time.Minute,
	}
}

type SessionSetInput struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	Client    string    `json:"client"`
	LastUsed  time.Time `json:"last_used,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (s Sessions) Set(ctx context.Context, input SessionCreateInput) error {
	sessionKey := fmt.Sprintf("%s:id:%s", sessionsCollection, input.ID.String())
	userSessionsKey := fmt.Sprintf("%s:user:%s", sessionsCollection, input.UserID.String())

	//// Проверка на существование сессии.
	//exists, err := s.client.Exists(ctx, sessionKey).Result()
	//if err != nil {
	//	return fmt.Errorf("error checking session existence: %w", err)
	//}
	//if exists > 0 {
	//	return fmt.Errorf("session already exists")
	//}

	if err := s.client.Del(ctx, sessionKey).Err(); err != nil {
		//return fmt.Errorf("error deleting existing session key: %w", err)
	}

	// Обновляем индекс email: в emailIndexKey ставим для поля input.Email значение input.ID.
	if err := s.client.HDel(ctx, userSessionsKey, input.UserID.String()).Err(); err != nil {
		// Если ключа нет, можно проигнорировать ошибку.
	}

	data := map[string]interface{}{
		"user_id":    input.UserID.String(),
		"token":      input.Token,
		"client":     input.Client,
		"created_at": input.CreatedAt.Format(time.RFC3339),
		"last_used":  input.LastUsed.Format(time.RFC3339),
	}

	if err := s.client.HSet(ctx, sessionKey, data).Err(); err != nil {
		return fmt.Errorf("error adding session to Redis: %w", err)
	}

	if err := s.client.SAdd(ctx, userSessionsKey, input.ID.String()).Err(); err != nil {
		return fmt.Errorf("error indexing session under user: %w", err)
	}

	if s.lifeTime > 0 {
		pipe := s.client.Pipeline()
		pipe.Expire(ctx, sessionKey, s.lifeTime)
		pipe.Expire(ctx, userSessionsKey, s.lifeTime)
		if _, err := pipe.Exec(ctx); err != nil {
			return fmt.Errorf("error setting expiration for session keys: %w", err)
		}
	}

	return nil
}

type SessionCreateInput struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	Client    string    `json:"client"`
	LastUsed  time.Time `json:"last_used,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Create создаёт новую сессию, если с таким sessionID еще нет.
// Помимо сохранения данных в hash по ключу sessions:id:<sessionID>,
// идентификатор сессии добавляется в множество sessions:user:<userID>.
func (s Sessions) Create(ctx context.Context, input SessionCreateInput) error {
	sessionKey := fmt.Sprintf("%s:id:%s", sessionsCollection, input.ID.String())
	userSessionsKey := fmt.Sprintf("%s:user:%s", sessionsCollection, input.UserID.String())

	// Проверка на существование сессии.
	exists, err := s.client.Exists(ctx, sessionKey).Result()
	if err != nil {
		return fmt.Errorf("error checking session existence: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("session already exists")
	}

	data := map[string]interface{}{
		"user_id":    input.UserID.String(),
		"token":      input.Token,
		"client":     input.Client,
		"created_at": input.CreatedAt.Format(time.RFC3339),
		"last_used":  input.LastUsed.Format(time.RFC3339),
	}

	if err := s.client.HSet(ctx, sessionKey, data).Err(); err != nil {
		return fmt.Errorf("error adding session to Redis: %w", err)
	}

	if err := s.client.SAdd(ctx, userSessionsKey, input.ID.String()).Err(); err != nil {
		return fmt.Errorf("error indexing session under user: %w", err)
	}

	if s.lifeTime > 0 {
		pipe := s.client.Pipeline()
		pipe.Expire(ctx, sessionKey, s.lifeTime)
		pipe.Expire(ctx, userSessionsKey, s.lifeTime)
		if _, err := pipe.Exec(ctx); err != nil {
			return fmt.Errorf("error setting expiration for session keys: %w", err)
		}
	}

	return nil
}

// GetByID возвращает данные сессии по sessionID.
func (s Sessions) GetByID(ctx context.Context, sessionID string) (SessionModel, error) {
	sessionKey := fmt.Sprintf("%s:id:%s", sessionsCollection, sessionID)

	vals, err := s.client.HGetAll(ctx, sessionKey).Result()
	if err != nil {
		return SessionModel{}, fmt.Errorf("error getting session data: %w", err)
	}
	if len(vals) == 0 {
		return SessionModel{}, fmt.Errorf("session not found, id=%s", sessionID)
	}
	return parseSession(sessionID, vals)
}

func (s Sessions) GetByUserID(ctx context.Context, userID string) ([]SessionModel, error) {
	userSessionsKey := fmt.Sprintf("%s:user:%s", sessionsCollection, userID)
	sessionIDs, err := s.client.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting sessions for user: %w", err)
	}

	var sessionsArr []SessionModel
	for _, id := range sessionIDs {
		sessionKey := fmt.Sprintf("%s:id:%s", sessionsCollection, id)
		vals, err := s.client.HGetAll(ctx, sessionKey).Result()
		if err != nil {
			return nil, fmt.Errorf("error getting session %s: %w", id, err)
		}
		ses, err := parseSession(id, vals)
		if err != nil {
			return nil, fmt.Errorf("error parsing session %s: %w", id, err)
		}
		sessionsArr = append(sessionsArr, ses)
	}
	return sessionsArr, nil
}

type SessionUpdateInput struct {
	Token    *string   `json:"token"`
	LastUsed time.Time `json:"last_used,omitempty"`
}

// Update обновляет данные сессии. Помимо обновления hash по ключу sessions:id:<sessionID>,
// методом также гарантируется, что сессия принадлежит нужному аккаунту.
// Затем обновление отражается и в индексе sessions:user:<userID> за счёт продления времени жизни.
func (s Sessions) Update(ctx context.Context, sessionID, userID uuid.UUID, update SessionUpdateInput) error {
	sessionKey := fmt.Sprintf("%s:id:%s", sessionsCollection, sessionID.String())

	ses, err := s.GetByID(ctx, sessionID.String())
	if err != nil {
		return fmt.Errorf("failed to get session details: %w", err)
	}
	if ses.UserID != userID {
		return fmt.Errorf("session %s does not belong to user %s", sessionID, userID)
	}

	data := make(map[string]interface{})
	if update.Token != nil {
		data["token"] = *update.Token
	}
	if !update.LastUsed.IsZero() {
		data["last_used"] = update.LastUsed.Format(time.RFC3339)
	}

	userSessionsKey := fmt.Sprintf("%s:user:%s", sessionsCollection, userID.String())

	pipe := s.client.Pipeline()
	pipe.HSet(ctx, sessionKey, data)
	// Добавляем sessionID в множество, чтобы убедиться, что оно там присутствует.
	pipe.SAdd(ctx, userSessionsKey, sessionID.String())
	pipe.Expire(ctx, sessionKey, s.lifeTime)
	pipe.Expire(ctx, userSessionsKey, s.lifeTime)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("error executing session update pipeline: %w", err)
	}

	return nil
}

// Delete удаляет сессию по sessionID и удаляет ее из множества сессий аккаунта.
func (s Sessions) Delete(ctx context.Context, sessionID uuid.UUID) error {
	// Получаем данные сессии, чтобы узнать userID.
	ses, err := s.GetByID(ctx, sessionID.String())
	if err != nil {
		return fmt.Errorf("failed to retrieve session: %w", err)
	}

	sessionKey := fmt.Sprintf("%s:id:%s", sessionsCollection, sessionID.String())
	userSessionsKey := fmt.Sprintf("%s:user:%s", sessionsCollection, ses.UserID.String())

	// Опционально: проверяем, что sessionID действительно присутствует в множестве.
	isMember, err := s.client.SIsMember(ctx, userSessionsKey, sessionID.String()).Result()
	if err != nil {
		return fmt.Errorf("error checking session membership: %w", err)
	}
	if !isMember {
		// Логгировать предупреждение, но продолжаем удаление.
	}

	pipe := s.client.Pipeline()
	pipe.Del(ctx, sessionKey)
	pipe.SRem(ctx, userSessionsKey, sessionID.String())
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("error deleting session keys from Redis: %w", err)
	}

	return nil
}

func (s Sessions) Drop(ctx context.Context) error {
	pattern := fmt.Sprintf("%s:*", sessionsCollection)
	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("error fetching keys with pattern %s: %w", pattern, err)
	}
	if len(keys) == 0 {
		return nil
	}
	if err := s.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete keys with pattern %s: %w", pattern, err)
	}
	return nil
}

func parseSession(sessionID string, vals map[string]string) (SessionModel, error) {
	createdAt, err := time.Parse(time.RFC3339, vals["created_at"])
	if err != nil {
		return SessionModel{}, fmt.Errorf("error parsing created_at: %w", err)
	}

	userID, err := uuid.Parse(vals["user_id"])
	if err != nil {
		return SessionModel{}, fmt.Errorf("error parsing user_id: %w", err)
	}

	sID, err := uuid.Parse(sessionID) // sessionID здесь передается как чистый идентификатор (без префикса)
	if err != nil {
		return SessionModel{}, fmt.Errorf("error parsing sessionID: %w", err)
	}

	res := SessionModel{
		ID:        sID,
		UserID:    userID,
		Token:     vals["token"],
		Client:    vals["client"],
		CreatedAt: createdAt,
	}

	if lastUsed, ok := vals["last_used"]; ok && lastUsed != "" {
		lu, err := time.Parse(time.RFC3339, lastUsed)
		if err != nil {
			return SessionModel{}, fmt.Errorf("error parsing last_used: %w", err)
		}
		res.LastUsed = lu
	}

	return res, nil
}

func (s Sessions) Terminate(ctx context.Context, userID uuid.UUID) error {
	// Формируем ключ для индекса сессий аккаунта.
	userSessionsKey := fmt.Sprintf("%s:user:%s", sessionsCollection, userID.String())

	// Получаем список sessionID, связанных с этим аккаунтом.
	sessionIDs, err := s.client.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return fmt.Errorf("error getting sessions for user %s: %w", userID.String(), err)
	}

	// Если сессий нет, можно сразу вернуть.
	if len(sessionIDs) == 0 {
		return nil
	}

	// Открываем пайплайн для группового удаления.
	pipe := s.client.Pipeline()
	for _, sid := range sessionIDs {
		// Для каждого sessionID формируем ключ хэша сессии.
		sessionKey := fmt.Sprintf("%s:id:%s", sessionsCollection, sid)
		pipe.Del(ctx, sessionKey)
	}

	// Удаляем индексный ключ для аккаунта.
	pipe.Del(ctx, userSessionsKey)

	// Выполняем пайплайн.
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("error deleting sessions for user %s: %w", userID.String(), err)
	}

	return nil
}
