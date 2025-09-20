package apptest

import (
	"context"
	"testing"

	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/google/uuid"
)

func CreateUser(s Setup, t *testing.T, email, password, role string) models.User {
	t.Helper()
	ctx := context.Background()

	user, err := s.app.Register_ONLY_FOR_TESTS(ctx, email, password, role)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	return user
}

func CreateSession(s Setup, t *testing.T, userID uuid.UUID) models.Session {
	t.Helper()
	ctx := context.Background()

	session, err := s.app.CreateSession_ONLY_FOR_TESTS(ctx, userID)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	return session
}
