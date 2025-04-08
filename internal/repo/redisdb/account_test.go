package redisdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/tokens/roles" // если roles.Role является типом alias для string
	"github.com/redis/go-redis/v9"
)

func setupAccounts(t *testing.T) (Accounts, func()) {
	// Создаем Redis клиент для тестовой базы (например, DB=1)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:7200",
		DB:   0,
	})
	accounts := NewAccounts(client, 5) // срок жизни 5 минут

	// Очистка: удалим все ключи, начинающиеся с "accounts:"
	ctx := context.Background()
	if err := accounts.Drop(ctx); err != nil {
		t.Fatalf("failed to drop account keys: %v", err)
	}

	cleanup := func() {
		client.Close()
	}

	return accounts, cleanup
}

func TestAccountsCreateAndGet(t *testing.T) {
	ctx := context.Background()
	accounts, cleanup := setupAccounts(t)
	defer cleanup()

	accountID := uuid.New()
	createdAt := time.Now().UTC()
	input := CreateAccountInput{
		ID:           accountID,
		Email:        "test@example.com",
		Role:         roles.Admin,
		Subscription: uuid.Nil,
		CreatedAt:    createdAt,
	}

	if err := accounts.Create(ctx, input); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := accounts.GetByID(ctx, accountID.String())
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Email != input.Email {
		t.Errorf("expected email %s, got %s", input.Email, got.Email)
	}
	if got.Role != string(input.Role) {
		t.Errorf("expected role %s, got %s", input.Role, got.Role)
	}
	if got.CreatedAt.Format(time.RFC3339) != createdAt.Format(time.RFC3339) {
		t.Errorf("expected created_at %v, got %v", createdAt, got.CreatedAt)
	}

	got2, err := accounts.GetByEmail(ctx, input.Email)
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}
	if got2.ID != accountID {
		t.Errorf("expected accountID %v, got %v", accountID, got2.ID)
	}
}

func TestAccountsSet(t *testing.T) {
	ctx := context.Background()
	accounts, cleanup := setupAccounts(t)
	defer cleanup()

	accountID := uuid.New()
	createdAt := time.Now().UTC()
	input := AccountSetInput{
		ID:           accountID,
		Email:        "set@example.com",
		Role:         roles.Role("admin"),
		Subscription: uuid.New(),
		CreatedAt:    createdAt,
	}
	if err := accounts.Set(ctx, input); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := accounts.GetByID(ctx, accountID.String())
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Email != input.Email {
		t.Errorf("expected email %s, got %s", input.Email, got.Email)
	}
	if got.Role != string(input.Role) {
		t.Errorf("expected role %s, got %s", input.Role, got.Role)
	}
}

func TestAccountsUpdate(t *testing.T) {
	ctx := context.Background()
	accounts, cleanup := setupAccounts(t)
	defer cleanup()

	accountID := uuid.New()
	createdAt := time.Now().UTC()
	// Сначала создаём аккаунт.
	createInput := CreateAccountInput{
		ID:           accountID,
		Email:        "update@example.com",
		Role:         roles.Role("user"),
		Subscription: uuid.Nil,
		CreatedAt:    createdAt,
	}
	if err := accounts.Create(ctx, createInput); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Подождем, чтобы время обновления отличалось.
	time.Sleep(1 * time.Second)
	updateTime := time.Now().UTC()
	newRole := roles.Role("super_user")
	newSub := uuid.New()

	updateReq := AccountUpdateRequest{
		Role:         &newRole,
		Subscription: &newSub,
		UpdatedAt:    updateTime,
	}
	if err := accounts.Update(ctx, accountID, updateReq); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, err := accounts.GetByID(ctx, accountID.String())
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Role != string(newRole) {
		t.Errorf("expected role %s, got %s", newRole, got.Role)
	}
	if got.Subscription != newSub {
		t.Errorf("expected subscription %s, got %s", newSub, got.Subscription)
	}
	if got.UpdatedAt == nil || got.UpdatedAt.Format(time.RFC3339) != updateTime.Format(time.RFC3339) {
		t.Errorf("expected updated_at %v, got %v", updateTime, got.UpdatedAt)
	}
}

func TestAccountsDelete(t *testing.T) {
	ctx := context.Background()
	accounts, cleanup := setupAccounts(t)
	defer cleanup()

	accountID := uuid.New()
	createdAt := time.Now().UTC()
	input := CreateAccountInput{
		ID:           accountID,
		Email:        "delete@example.com",
		Role:         roles.Role("user"),
		Subscription: uuid.Nil,
		CreatedAt:    createdAt,
	}
	if err := accounts.Create(ctx, input); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Убедимся, что аккаунт существует.
	_, err := accounts.GetByID(ctx, accountID.String())
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	// Удаляем аккаунт.
	if err := accounts.Delete(ctx, accountID.String(), input.Email); err != nil && err.Error() != "redis: nil" {
		t.Fatalf("Delete failed: %v", err)
	}

	// Проверяем, что аккаунт больше не существует.
	_, err = accounts.GetByID(ctx, accountID.String())
	if err == nil {
		t.Fatalf("expected error when getting deleted account, got nil")
	}
}

func TestAccountsDrop(t *testing.T) {
	ctx := context.Background()
	accounts, cleanup := setupAccounts(t)
	defer cleanup()

	// Создаем несколько аккаунтов.
	for i := 0; i < 3; i++ {
		accountID := uuid.New()
		email := fmt.Sprintf("drop%d@example.com", i)
		input := CreateAccountInput{
			ID:           accountID,
			Email:        email,
			Role:         roles.Role("user"),
			Subscription: uuid.Nil,
			CreatedAt:    time.Now().UTC(),
		}
		if err := accounts.Create(ctx, input); err != nil {
			t.Fatalf("Create failed for account %s: %v", accountID, err)
		}
	}

	// Проверяем, что какие-то ключи существуют.
	pattern := fmt.Sprintf("%s:*", accountsCollection)
	keys, err := accounts.client.Keys(ctx, pattern).Result()
	if err != nil {
		t.Fatalf("Keys command failed: %v", err)
	}
	if len(keys) == 0 {
		t.Fatalf("expected keys to exist before Drop")
	}

	// Вызываем Drop.
	if err := accounts.Drop(ctx); err != nil {
		t.Fatalf("Drop failed: %v", err)
	}

	// Проверяем, что ключи удалены.
	keys, err = accounts.client.Keys(ctx, pattern).Result()
	if err != nil {
		t.Fatalf("Keys command failed after Drop: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("expected 0 keys after Drop, got %d", len(keys))
	}
}
