package apptest

import (
	"context"
	"testing"

	"github.com/chains-lab/sso-svc/internal/models"
	"github.com/google/uuid"
)

func CreateUser(s Setup, t *testing.T, email, password, role string) models.User {
	t.Helper()
	ctx := context.Background()

	u, err := s.app.User().Register(ctx, email, password, role)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	return u
}

func CreateSession(s Setup, t *testing.T, userID uuid.UUID) models.Session {
	t.Helper()
	ctx := context.Background()

	u, err := s.app.User().GetByID(ctx, userID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}

	tkn, err := s.app.Session().Create(ctx, u.ID, u.Role)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	session, err := s.app.Session().Get(ctx, tkn.SessionID)
	if err != nil {
		t.Fatalf("GetByAccessToken: %v", err)
	}

	return session
}
