package sqldb_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/chains-lab/chains-auth/internal/repo/sqldb"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:7000/postgres?sslmode=disable")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}
	return db
}

func cleanupTables(t *testing.T, db *sql.DB) {
	tables := []string{"users", "sessions"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			t.Fatalf("failed to clean table %s: %v", table, err)
		}
	}
}

func TestIntegration_UserAndSession(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	cleanupTables(t, db)

	ctx := context.Background()

	// Создаем объекты для работы с аккаунтами и сессиями.
	userQ := sqldb.NewUsers(db)
	sessionsQ := sqldb.NewSessions(db)

	// 1. Создаем аккаунт.
	accID := uuid.New()
	accCreatedAt := time.Now().UTC()
	accInput := sqldb.UserInsertInput{
		ID:           accID,
		Email:        "integration@example.com",
		Role:         roles.Role("user"),
		Subscription: uuid.Nil,
		CreatedAt:    accCreatedAt,
	}

	if err := userQ.Insert(ctx, accInput); err != nil {
		t.Fatalf("User Insert failed: %v", err)
	}

	// Проверяем создание аккаунта.
	user, err := userQ.FilterID(accID).Get(ctx)
	if err != nil {
		t.Fatalf("User GetByID failed: %v", err)
	}
	if user.Email != accInput.Email {
		t.Errorf("expected email %s, got %s", accInput.Email, user.Email)
	}

	// 2. Создаем сессию для этого аккаунта.
	sessID := uuid.New()
	sessCreatedAt := time.Now().UTC()
	lastUsed := sessCreatedAt.Add(1 * time.Minute)
	sessionInput := sqldb.SessionInsertInput{
		ID:        sessID,
		UserID:    accID,
		Token:     "test-token",
		Client:    "web",
		CreatedAt: sessCreatedAt,
		LastUsed:  lastUsed,
	}

	if err := sessionsQ.Insert(ctx, sessionInput); err != nil {
		t.Fatalf("Session Insert failed: %v", err)
	}

	// Проверяем получение сессии по ID.
	session, err := sessionsQ.FilterID(sessID).Get(ctx)
	if err != nil {
		t.Fatalf("Session GetByID failed: %v", err)
	}
	if session.Token != sessionInput.Token {
		t.Errorf("expected token %s, got %s", sessionInput.Token, session.Token)
	}
	//if !session.CreatedAt.Equal(sessCreatedAt) {
	//	t.Errorf("expected created_at %v, got %v", sessCreatedAt, session.CreatedAt)
	//}

	// Проверяем получение сессии по UserID.
	sessions, err := sessionsQ.FilterUserID(accID).Select(ctx)
	if err != nil {
		t.Fatalf("Sessions Select by UserID failed: %v", err)
	}
	if len(sessions) != 1 {
		t.Errorf("expected 1 session for user, got %d", len(sessions))
	}

	// 3. Обновляем аккаунт.
	newRole := roles.Role("admin")
	updateTime := time.Now().UTC()
	accUpdate := sqldb.UserUpdateInput{
		Role:         &newRole,
		Subscription: nil, // оставляем без изменений
		UpdatedAt:    updateTime,
	}
	// Применяем фильтр по ID.
	if err := userQ.FilterID(accID).Update(ctx, accUpdate); err != nil {
		t.Fatalf("User Update failed: %v", err)
	}

	updatedUser, err := userQ.FilterID(accID).Get(ctx)
	if err != nil {
		t.Fatalf("User GetByID after update failed: %v", err)
	}
	if updatedUser.Role != newRole {
		t.Errorf("expected updated role %s, got %s", newRole, updatedUser.Role)
	}
	// Если поле updated_at не было установлено ранее, оно должно быть равно updateTime.
	//if updatedUser.UpdatedAt == nil || !updatedUser.UpdatedAt.Equal(updateTime) {
	//	t.Errorf("expected updated_at %v, got %v", updateTime, updatedUser.UpdatedAt)
	//}

	// 4. Обновляем сессию.
	newToken := "updated-token"
	sessUpdate := sqldb.SessionUpdateInput{
		Token:    &newToken,
		LastUsed: time.Now().UTC(),
	}
	// Применяем фильтр по sessionID.
	if err := sessionsQ.FilterID(sessID).Update(ctx, sessUpdate); err != nil {
		t.Fatalf("Session Update failed: %v", err)
	}

	updatedSession, err := sessionsQ.FilterID(sessID).Get(ctx)
	if err != nil {
		t.Fatalf("Session GetByID after update failed: %v", err)
	}
	if updatedSession.Token != newToken {
		t.Errorf("expected updated token %s, got %s", newToken, updatedSession.Token)
	}

	// 5. Удаляем сессию.
	if err := sessionsQ.FilterID(sessID).Delete(ctx); err != nil {
		t.Fatalf("Session Delete failed: %v", err)
	}
	_, err = sessionsQ.FilterID(sessID).Get(ctx)
	if err == nil {
		t.Errorf("expected error fetching deleted session, got none")
	}

	// 6. Удаляем аккаунт.
	if err := userQ.FilterID(accID).Delete(ctx); err != nil {
		t.Fatalf("User Delete failed: %v", err)
	}
	_, err = userQ.FilterID(accID).Get(ctx)
	if err == nil {
		t.Errorf("expected error fetching deleted user, got none")
	}
}
