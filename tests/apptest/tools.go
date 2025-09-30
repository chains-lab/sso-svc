package apptest

import (
	"context"
	"testing"

	models2 "github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

func CreateUser(s Setup, t *testing.T, email, password, role string) models2.User {
	t.Helper()
	ctx := context.Background()

	u, err := s.core.Auth.Register(ctx, email, password, role)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	return u
}

func CreateSession(s Setup, t *testing.T, userID uuid.UUID) models2.Session {
	t.Helper()
	ctx := context.Background()

	u, err := s.core.User.GetByID(ctx, userID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}

	tkn, err := s.core.Auth.CreateSession(ctx, u.ID, u.Role)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	session, err := s.core.Session.Get(ctx, tkn.SessionID)
	if err != nil {
		t.Fatalf("GetByAccessToken: %v", err)
	}

	return session
}
