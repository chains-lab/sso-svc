package apptest

import (
	"context"
	"testing"

	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/chains-lab/sso-svc/internal/domain/services/user"
	"github.com/google/uuid"
)

func CreateUser(s Setup, t *testing.T, email, password, role string) models.User {
	t.Helper()
	ctx := context.Background()

	u, err := s.domain.user.Register(ctx, user.RegisterParams{
		Email:    email,
		Password: password,
		Role:     role,
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	return u
}

func CreateSession(s Setup, t *testing.T, userID uuid.UUID) models.Session {
	t.Helper()
	ctx := context.Background()

	u, err := s.domain.user.GetByID(ctx, userID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}

	tkn, err := s.domain.session.Create(ctx, u.ID, u.Role, u.EmailVer)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	session, err := s.domain.session.Get(ctx, tkn.SessionID)
	if err != nil {
		t.Fatalf("GetByAccessToken: %v", err)
	}

	return session
}
