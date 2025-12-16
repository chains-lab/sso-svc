package domain_test

import (
	"context"
	"testing"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/google/uuid"
)

func CreateUser(s Setup, t *testing.T, email, password, role string) entity.User {
	t.Helper()
	ctx := context.Background()

	u, err := s.core.Auth.Register(ctx, email, password, role)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	return u
}

func CreateSession(s Setup, t *testing.T, userID uuid.UUID) entity.Session {
	t.Helper()
	ctx := context.Background()

	u, err := s.core.User.GetByID(ctx, userID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}

	tkn, err := s.core.Auth.CreateSession(ctx, u)
	if err != nil {
		t.Fatalf("createSession: %v", err)
	}

	session, err := s.core.Session.Get(ctx, tkn.SessionID)
	if err != nil {
		t.Fatalf("GetByAccessToken: %v", err)
	}

	return session
}
