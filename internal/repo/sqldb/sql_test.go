package sqldb_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/repo/sqldb"
	"github.com/hs-zavet/tokens/roles"
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
	tables := []string{"accounts", "sessions"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			t.Fatalf("failed to clean table %s: %v", table, err)
		}
	}
}

func TestIntegration_AccountAndSession(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	cleanupTables(t, db)

	ctx := context.Background()

	// Создаем объекты для работы с аккаунтами и сессиями.
	accountQ := sqldb.NewAccounts(db)
	sessionsQ := sqldb.NewSessions(db)

	// 1. Создаем аккаунт.
	accID := uuid.New()
	accCreatedAt := time.Now().UTC()
	accInput := sqldb.AccountInsertInput{
		ID:           accID,
		Email:        "integration@example.com",
		Role:         roles.Role("user"),
		Subscription: uuid.Nil,
		CreatedAt:    accCreatedAt,
	}

	if err := accountQ.Insert(ctx, accInput); err != nil {
		t.Fatalf("Account Insert failed: %v", err)
	}

	// Проверяем создание аккаунта.
	account, err := accountQ.FilterID(accID).Get(ctx)
	if err != nil {
		t.Fatalf("Account GetByID failed: %v", err)
	}
	if account.Email != accInput.Email {
		t.Errorf("expected email %s, got %s", accInput.Email, account.Email)
	}

	// 2. Создаем сессию для этого аккаунта.
	sessID := uuid.New()
	sessCreatedAt := time.Now().UTC()
	lastUsed := sessCreatedAt.Add(1 * time.Minute)
	sessionInput := sqldb.SessionInsertInput{
		ID:        sessID,
		AccountID: accID,
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

	// Проверяем получение сессии по AccountID.
	sessions, err := sessionsQ.FilterAccountID(accID).Select(ctx)
	if err != nil {
		t.Fatalf("Sessions Select by AccountID failed: %v", err)
	}
	if len(sessions) != 1 {
		t.Errorf("expected 1 session for account, got %d", len(sessions))
	}

	// 3. Обновляем аккаунт.
	newRole := roles.Role("admin")
	updateTime := time.Now().UTC()
	accUpdate := sqldb.AccountUpdateInput{
		Role:         &newRole,
		Subscription: nil, // оставляем без изменений
		UpdatedAt:    updateTime,
	}
	// Применяем фильтр по ID.
	if err := accountQ.FilterID(accID).Update(ctx, accUpdate); err != nil {
		t.Fatalf("Account Update failed: %v", err)
	}

	updatedAccount, err := accountQ.FilterID(accID).Get(ctx)
	if err != nil {
		t.Fatalf("Account GetByID after update failed: %v", err)
	}
	if updatedAccount.Role != newRole {
		t.Errorf("expected updated role %s, got %s", newRole, updatedAccount.Role)
	}
	// Если поле updated_at не было установлено ранее, оно должно быть равно updateTime.
	//if updatedAccount.UpdatedAt == nil || !updatedAccount.UpdatedAt.Equal(updateTime) {
	//	t.Errorf("expected updated_at %v, got %v", updateTime, updatedAccount.UpdatedAt)
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
	if err := accountQ.FilterID(accID).Delete(ctx); err != nil {
		t.Fatalf("Account Delete failed: %v", err)
	}
	_, err = accountQ.FilterID(accID).Get(ctx)
	if err == nil {
		t.Errorf("expected error fetching deleted account, got none")
	}
}
