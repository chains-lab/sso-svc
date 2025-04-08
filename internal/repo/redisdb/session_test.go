package redisdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// setupSessions создаёт Redis-клиент для тестовой базы (например, DB=2)
// и возвращает инстанс Sessions с указанным временем жизни (например, 5 минут).
// После тестов вызывается cleanup для закрытия клиента.
func setupSessions(t *testing.T) (Sessions, func()) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0, // используем тестовую базу, чтобы не портить продакшн-данные
	})
	sess := NewSessions(client, 5) // 5 минут времени жизни

	// Очистим все ключи, начинающиеся с "sessions:"
	ctx := context.Background()
	if err := sess.Drop(ctx); err != nil {
		t.Fatalf("failed to drop session keys: %v", err)
	}

	cleanup := func() {
		client.Close()
	}
	return sess, cleanup
}

func TestSessions_CreateAndGetByID(t *testing.T) {
	ctx := context.Background()
	sess, cleanup := setupSessions(t)
	defer cleanup()

	id := uuid.New()
	accountID := uuid.New()
	createdAt := time.Now().UTC()
	lastUsed := createdAt.Add(1 * time.Minute)

	input := SessionCreateInput{
		ID:        id,
		AccountID: accountID,
		Token:     "test-token",
		Client:    "test-client",
		CreatedAt: createdAt,
		LastUsed:  lastUsed,
	}

	if err := sess.Create(ctx, input); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := sess.GetByID(ctx, id.String())
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.ID != id {
		t.Errorf("expected session id %v, got %v", id, got.ID)
	}
	if got.AccountID != accountID {
		t.Errorf("expected account id %v, got %v", accountID, got.AccountID)
	}
	if got.Token != input.Token {
		t.Errorf("expected token %s, got %s", input.Token, got.Token)
	}
	if got.Client != input.Client {
		t.Errorf("expected client %s, got %s", input.Client, got.Client)
	}
}

func TestSessions_GetByAccountID(t *testing.T) {
	ctx := context.Background()
	sess, cleanup := setupSessions(t)
	defer cleanup()

	accountID := uuid.New()
	// Создаем две сессии для одного аккаунта.
	sessionIDs := []uuid.UUID{uuid.New(), uuid.New()}
	for _, sid := range sessionIDs {
		input := SessionCreateInput{
			ID:        sid,
			AccountID: accountID,
			Token:     fmt.Sprintf("token-%s", sid.String()[:8]),
			Client:    "test-client",
			CreatedAt: time.Now().UTC(),
			LastUsed:  time.Now().UTC(),
		}
		if err := sess.Create(ctx, input); err != nil {
			t.Fatalf("failed to create session %s: %v", sid, err)
		}
	}

	sessions, err := sess.GetByAccountID(ctx, accountID.String())
	if err != nil {
		t.Fatalf("GetByAccountID failed: %v", err)
	}
	if len(sessions) != len(sessionIDs) {
		t.Errorf("expected %d sessions, got %d", len(sessionIDs), len(sessions))
	}
	// Можно проверить, что все созданные sessionID присутствуют.
	foundCount := 0
	for _, s := range sessions {
		for _, sid := range sessionIDs {
			if s.ID == sid {
				foundCount++
			}
		}
	}
	if foundCount != len(sessionIDs) {
		t.Errorf("not all sessions found; expected %d, got %d", len(sessionIDs), foundCount)
	}
}

func TestSessions_Update(t *testing.T) {
	ctx := context.Background()
	sess, cleanup := setupSessions(t)
	defer cleanup()

	id := uuid.New()
	accountID := uuid.New()
	createdAt := time.Now().UTC()
	lastUsed := createdAt.Add(1 * time.Minute)

	// Создаем первоначальную сессию.
	input := SessionCreateInput{
		ID:        id,
		AccountID: accountID,
		Token:     "initial-token",
		Client:    "test-client",
		CreatedAt: createdAt,
		LastUsed:  lastUsed,
	}
	if err := sess.Create(ctx, input); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Формируем данные для обновления.
	newToken := "updated-token"
	newLastUsed := time.Now().UTC().Add(2 * time.Minute)
	updateInput := SessionUpdateInput{
		Token:    &newToken,
		LastUsed: newLastUsed,
	}

	// Обновляем сессию; также проверяется принадлежность: мы передаем accountID.
	if err := sess.Update(ctx, id, accountID, updateInput); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Получаем обновленные данные.
	updated, err := sess.GetByID(ctx, id.String())
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if updated.Token != newToken {
		t.Errorf("expected token %s, got %s", newToken, updated.Token)
	}
	// Проверка времени — сравниваем по формату.
	if updated.LastUsed.Format(time.RFC3339) != newLastUsed.Format(time.RFC3339) {
		t.Errorf("expected last_used %v, got %v", newLastUsed, updated.LastUsed)
	}
}

func TestSessions_Delete(t *testing.T) {
	ctx := context.Background()
	sess, cleanup := setupSessions(t)
	defer cleanup()

	id := uuid.New()
	accountID := uuid.New()
	createdAt := time.Now().UTC()
	lastUsed := createdAt.Add(1 * time.Minute)

	input := SessionCreateInput{
		ID:        id,
		AccountID: accountID,
		Token:     "token-to-delete",
		Client:    "test-client",
		CreatedAt: createdAt,
		LastUsed:  lastUsed,
	}
	if err := sess.Create(ctx, input); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Удаляем сессию.
	if err := sess.Delete(ctx, id); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Проверяем, что GetByID возвращает ошибку.
	_, err := sess.GetByID(ctx, id.String())
	if err == nil {
		t.Fatalf("expected error on GetByID after deletion, got nil")
	}
}

func TestSessions_Terminate(t *testing.T) {
	ctx := context.Background()
	sess, cleanup := setupSessions(t)
	defer cleanup()

	accountID := uuid.New()

	// Создаем несколько сессий для одного аккаунта.
	var sessionIDs []uuid.UUID
	for i := 0; i < 3; i++ {
		sid := uuid.New()
		sessionIDs = append(sessionIDs, sid)
		input := SessionCreateInput{
			ID:        sid,
			AccountID: accountID,
			Token:     fmt.Sprintf("token-%d", i),
			Client:    "test-client",
			CreatedAt: time.Now().UTC(),
			LastUsed:  time.Now().UTC(),
		}
		if err := sess.Create(ctx, input); err != nil {
			t.Fatalf("Create failed for session %d: %v", i, err)
		}
	}

	// Убедимся, что сессии существуют.
	sessionsBefore, err := sess.GetByAccountID(ctx, accountID.String())
	if err != nil {
		t.Fatalf("GetByAccountID failed: %v", err)
	}
	if len(sessionsBefore) != len(sessionIDs) {
		t.Fatalf("expected %d sessions before termination, got %d", len(sessionIDs), len(sessionsBefore))
	}

	// Вызываем Terminate — удаляем все сессии для аккаунта.
	if err := sess.Terminate(ctx, accountID); err != nil {
		t.Fatalf("Terminate failed: %v", err)
	}

	// Проверяем, что сессии для аккаунта отсутствуют.
	sessionsAfter, err := sess.GetByAccountID(ctx, accountID.String())
	if err != nil {
		t.Fatalf("GetByAccountID failed after termination: %v", err)
	}
	if len(sessionsAfter) != 0 {
		t.Errorf("expected 0 sessions after termination, got %d", len(sessionsAfter))
	}
}
