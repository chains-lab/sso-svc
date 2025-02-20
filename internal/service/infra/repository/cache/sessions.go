package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/models"
	"github.com/redis/go-redis/v9"
)

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

func (s *Sessions) Add(ctx context.Context, session models.Session) error {
	sessionKey := fmt.Sprintf("session:id:%s", session.ID)
	userSessionsKey := fmt.Sprintf("session:user:%s", session.UserID)

	exists, err := s.client.Exists(ctx, sessionKey).Result()
	if err != nil {
		return fmt.Errorf("error checking session existence: %w", err)
	}

	if exists > 0 {
		updateData := map[string]interface{}{
			"last_used": session.LastUsed.Format(time.RFC3339),
			"token":     session.Token,
		}

		err = s.client.HSet(ctx, sessionKey, updateData).Err()
		if err != nil {
			return fmt.Errorf("error updating existing session: %w", err)
		}
	} else {
		data := map[string]interface{}{
			"user_id":   session.UserID.String(),
			"token":     session.Token,
			"client":    session.Client,
			"ip":        session.IP,
			"create_at": session.CreatedAt.Format(time.RFC3339),
			"last_used": session.LastUsed.Format(time.RFC3339),
		}

		err = s.client.HSet(ctx, sessionKey, data).Err()
		if err != nil {
			return fmt.Errorf("error adding session to Redis: %w", err)
		}

		err = s.client.SAdd(ctx, userSessionsKey, session.ID.String()).Err()
		if err != nil {
			return fmt.Errorf("error indexing session under user_id: %w", err)
		}
	}

	if s.lifeTime > 0 {
		_ = s.client.Expire(ctx, sessionKey, s.lifeTime).Err()
		_ = s.client.Expire(ctx, userSessionsKey, s.lifeTime).Err()
	}

	return nil
}

func (s *Sessions) GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	key := fmt.Sprintf("session:id:%s", sessionID)

	vals, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting session from Redis: %w", err)
	}

	if len(vals) == 0 {
		return nil, fmt.Errorf("sesion not found, id=%s", sessionID)
	}

	createdAt, err := time.Parse(time.RFC3339, vals["create_at"])
	if err != nil {
		return nil, fmt.Errorf("error parsing create_at: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, vals["last_used"])
	if err != nil {
		return nil, fmt.Errorf("error parsing updated_at: %w", err)
	}

	userID, err := uuid.Parse(vals["user_id"])
	if err != nil {
		return nil, fmt.Errorf("error parsing user ID: %w", err)
	}

	session := &models.Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     vals["token"],
		Client:    vals["client"],
		IP:        vals["ip"],
		LastUsed:  updatedAt,
		CreatedAt: createdAt,
	}

	return session, nil
}

func (s *Sessions) SelectByUserID(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	userSessionsKey := fmt.Sprintf("session:user:%s", userID)

	sessionIDs, err := s.client.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting Sessions by user_id: %w", err)
	}

	var sessionsArr []models.Session
	for _, sessionID := range sessionIDs {
		id, err := uuid.Parse(sessionID)
		if err != nil {
			return nil, fmt.Errorf("error parsing session ID: %w", err)
		}
		session, err := s.GetByID(ctx, id)
		if err == nil {
			sessionsArr = append(sessionsArr, *session)
		}
	}

	return sessionsArr, nil
}

func (s *Sessions) DeleteByUserID(ctx context.Context, userID uuid.UUID, curSessionID *uuid.UUID) error {
	userSessionsKey := fmt.Sprintf("session:user:%s", userID)

	sessionIDs, err := s.client.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return fmt.Errorf("error getting Sessions for user: %w", err)
	}

	for _, sessionID := range sessionIDs {
		id, err := uuid.Parse(sessionID)
		if err != nil {
			return fmt.Errorf("error parsing session ID: %w", err)
		}
		if id != *curSessionID {
			_ = s.Delete(ctx, id)
		}
	}

	err = s.client.Del(ctx, userSessionsKey).Err()
	if err != nil {
		return fmt.Errorf("error deleting session list for user: %w", err)
	}

	return nil
}

func (s *Sessions) Delete(ctx context.Context, sessionID uuid.UUID) error {
	key := fmt.Sprintf("session:id:%s", sessionID)

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

	return nil
}
