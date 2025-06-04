package redisdb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chains-lab/chains-auth/internal/repo/redisdb"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func setupUsers(t *testing.T) (redisdb.Users, func()) {
	// Создаем Redis клиент для тестовой базы (например, DB=1)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:7200",
		DB:   0,
	})
	users := redisdb.NewUsers(client, 5) // срок жизни 5 минут

	// Очистка: удалим все ключи, начинающиеся с "users:"
	ctx := context.Background()
	if err := users.Drop(ctx); err != nil {
		t.Fatalf("failed to drop user keys: %v", err)
	}

	cleanup := func() {
		client.Close()
	}

	return users, cleanup
}

func TestUsersCreateAndGet(t *testing.T) {
	ctx := context.Background()
	users, cleanup := setupUsers(t)
	defer cleanup()

	userID := uuid.New()
	createdAt := time.Now().UTC()
	input := redisdb.CreateUserInput{
		ID:           userID,
		Email:        "test@example.com",
		Role:         roles.Admin,
		Subscription: uuid.Nil,
		CreatedAt:    createdAt,
	}

	if err := users.Create(ctx, input); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := users.GetByID(ctx, userID.String())
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

	got2, err := users.GetByEmail(ctx, input.Email)
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}
	if got2.ID != userID {
		t.Errorf("expected userID %v, got %v", userID, got2.ID)
	}
}

func TestUsersSet(t *testing.T) {
	ctx := context.Background()
	users, cleanup := setupUsers(t)
	defer cleanup()

	userID := uuid.New()
	createdAt := time.Now().UTC()
	input := redisdb.UserSetInput{
		ID:           userID,
		Email:        "set@example.com",
		Role:         roles.Role("admin"),
		Subscription: uuid.New(),
		CreatedAt:    createdAt,
	}
	if err := users.Set(ctx, input); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := users.GetByID(ctx, userID.String())
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

func TestUsersUpdate(t *testing.T) {
	ctx := context.Background()
	users, cleanup := setupUsers(t)
	defer cleanup()

	userID := uuid.New()
	createdAt := time.Now().UTC()
	// Сначала создаём аккаунт.
	createInput := redisdb.CreateUserInput{
		ID:           userID,
		Email:        "update@example.com",
		Role:         roles.Role("user"),
		Subscription: uuid.Nil,
		CreatedAt:    createdAt,
	}
	if err := users.Create(ctx, createInput); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Подождем, чтобы время обновления отличалось.
	time.Sleep(1 * time.Second)
	updateTime := time.Now().UTC()
	newRole := roles.Role("super_user")
	newSub := uuid.New()

	updateReq := redisdb.UserUpdateRequest{
		Role:         &newRole,
		Subscription: &newSub,
		UpdatedAt:    updateTime,
	}
	if err := users.Update(ctx, userID, updateReq); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, err := users.GetByID(ctx, userID.String())
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

func TestUsersDelete(t *testing.T) {
	ctx := context.Background()
	users, cleanup := setupUsers(t)
	defer cleanup()

	userID := uuid.New()
	createdAt := time.Now().UTC()
	input := redisdb.CreateUserInput{
		ID:           userID,
		Email:        "delete@example.com",
		Role:         roles.Role("user"),
		Subscription: uuid.Nil,
		CreatedAt:    createdAt,
	}
	if err := users.Create(ctx, input); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Убедимся, что аккаунт существует.
	_, err := users.GetByID(ctx, userID.String())
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	// Удаляем аккаунт.
	if err := users.Delete(ctx, userID.String(), input.Email); err != nil && err.Error() != "redis: nil" {
		t.Fatalf("Delete failed: %v", err)
	}

	// Проверяем, что аккаунт больше не существует.
	_, err = users.GetByID(ctx, userID.String())
	if err == nil {
		t.Fatalf("expected error when getting deleted user, got nil")
	}
}

func TestUsersDrop(t *testing.T) {
	ctx := context.Background()
	users, cleanup := setupUsers(t)
	defer cleanup()

	// Создаем несколько аккаунтов.
	var userIDs []uuid.UUID
	for i := 0; i < 3; i++ {
		aid := uuid.New()
		userIDs = append(userIDs, aid)
		email := fmt.Sprintf("drop%d@example.com", i)
		input := redisdb.CreateUserInput{
			ID:           aid,
			Email:        email,
			Role:         roles.Role("user"),
			Subscription: uuid.Nil,
			CreatedAt:    time.Now().UTC(),
		}
		if err := users.Create(ctx, input); err != nil {
			t.Fatalf("Create failed for user %s: %v", aid, err)
		}
	}

	// Убедимся, что аккаунты существуют.
	for _, aid := range userIDs {
		_, err := users.GetByID(ctx, aid.String())
		if err != nil {
			t.Fatalf("GetByID failed for user %s: %v", aid, err)
		}
	}

	// Вызываем Drop.
	if err := users.Drop(ctx); err != nil {
		t.Fatalf("Drop failed: %v", err)
	}
}
