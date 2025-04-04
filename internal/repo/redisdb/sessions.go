package redisdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const sessionsCollection = "sessions"

type SessionModel struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	Token     string    `json:"token"`
	Client    string    `json:"client"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}

type Sessions struct {
	client   *redis.Client
	lifeTime time.Duration
}

func NewSessions(client *redis.Client, lifeTime time.Duration) Sessions {
	return Sessions{
		client:   client,
		lifeTime: lifeTime,
	}
}

type AddSessionInput struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	Token     string    `json:"token"`
	Client    string    `json:"client"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}

func (s Sessions) Add(ctx context.Context, input AddSessionInput) error {
	sessionKey := fmt.Sprintf("%s:id:%s", sessionsCollection, input.ID)
	accountSessionsKey := fmt.Sprintf("%s:account:%s", sessionsCollection, input.AccountID)

	data := map[string]interface{}{
		"account_id": input.AccountID.String(),
		"token":      input.Token,
		"client":     input.Client,
		"ip":         input.IP,
		"create_at":  input.CreatedAt.Format(time.RFC3339),
		"last_used":  input.LastUsed.Format(time.RFC3339),
	}

	if err := s.client.HSet(ctx, sessionKey, data).Err(); err != nil {
		return fmt.Errorf("error adding input to Redis: %w", err)
	}

	if err := s.client.SAdd(ctx, accountSessionsKey, input.ID.String()).Err(); err != nil {
		return fmt.Errorf("error indexing input under account_id: %w", err)
	}

	if s.lifeTime > 0 {
		pipe := s.client.Pipeline()
		keys := []string{sessionKey, accountSessionsKey}
		for _, key := range keys {
			pipe.Expire(ctx, key, s.lifeTime)
		}
		_, err := pipe.Exec(ctx)
		if err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf("error setting expiration for keys: %w", err)
		}
	}

	return nil
}

func (s Sessions) GetByID(ctx context.Context, sessionID string) (SessionModel, error) {
	key := fmt.Sprintf("%s:id:%s", sessionsCollection, sessionID)

	vals, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return SessionModel{}, err
	}

	if len(vals) == 0 {
		return SessionModel{}, fmt.Errorf("sesion not found, id=%s", sessionID)
	}

	return parseSession(key, vals)
}

func (s Sessions) SelectByAccountID(ctx context.Context, AccountID string) ([]SessionModel, error) {
	accountSessionsKey := fmt.Sprintf("%s:account:%s", sessionsCollection, AccountID)

	sessionIDs, err := s.client.SMembers(ctx, accountSessionsKey).Result()
	if err != nil {
		return nil, err
	}

	var sessionsArr []SessionModel
	for _, ID := range sessionIDs {
		vals, err := s.client.HGetAll(ctx, ID).Result()
		if err != nil {
			return nil, fmt.Errorf("error geting session: %w", err)
		}
		ses, err := parseSession(ID, vals)
		if err != nil {
			return nil, err
		}
		sessionsArr = append(sessionsArr, ses)
	}

	return sessionsArr, nil
}

func (s Sessions) DeleteAllByAccountID(ctx context.Context, AccountID uuid.UUID, curSessionID *uuid.UUID) error {
	accountSessionsKey := fmt.Sprintf("%s:account:%s", sessionsCollection, AccountID)

	sessionIDs, err := s.client.SMembers(ctx, accountSessionsKey).Result()
	if err != nil {
		return fmt.Errorf("error getting sessions for account: %w", err)
	}

	for _, sessionID := range sessionIDs {
		id, err := uuid.Parse(sessionID)
		if err != nil {
			return fmt.Errorf("error parsing session ID: %w", err)
		}
		if curSessionID != nil {
			if id != *curSessionID {
				_ = s.Delete(ctx, id)
			}
		}
	}

	err = s.client.Del(ctx, accountSessionsKey).Err()
	if err != nil {
		return fmt.Errorf("error deleting session list for account: %w", err)
	}

	return nil
}

func (s Sessions) Delete(ctx context.Context, sessionID uuid.UUID) error {
	key := fmt.Sprintf("%s:id:%s", sessionsCollection, sessionID)

	ses, err := s.GetByID(ctx, key)
	if err != nil {
		return err
	}

	accountSessionsKey := fmt.Sprintf("%s:account:%s", sessionsCollection, ses.AccountID)

	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error checking session existence in Redis: %w", err)
	}

	if exists == 0 {
		return nil
	}

	err = s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting session from Redis: %w", err)
	}

	err = s.client.SRem(ctx, accountSessionsKey, sessionID.String()).Err()
	if err != nil {
		return fmt.Errorf("error removing session ID from account set: %w", err)
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
	CreatedAt, err := time.Parse(time.RFC3339, vals["created_at"])
	if err != nil {
		return SessionModel{}, fmt.Errorf("error parsing created_at: %w", err)
	}

	LastUsed, err := time.Parse(time.RFC3339, vals["last_used"])
	if err != nil {
		return SessionModel{}, fmt.Errorf("error parsing last_used: %w", err)
	}

	aID, err := uuid.Parse(vals["account_id"])
	if err != nil {
		return SessionModel{}, fmt.Errorf("error parsing account_id: %w", err)
	}

	sID, err := uuid.Parse(sessionID)
	if err != nil {
		return SessionModel{}, fmt.Errorf("error parsing sessionID: %w", err)
	}

	return SessionModel{
		ID:        sID,
		AccountID: aID,
		Token:     vals["token"],
		Client:    vals["client"],
		IP:        vals["ip"],
		CreatedAt: CreatedAt,
		LastUsed:  LastUsed,
	}, nil
}
