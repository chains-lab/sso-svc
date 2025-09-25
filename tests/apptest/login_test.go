package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
)

func TestUserRegistration(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	firstEmail := "tests@example"
	password := "Test@1234"

	_, err = s.domain.user.Register(ctx,
		firstEmail,
		password,
		roles.User,
	)
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	_, err = s.domain.session.Login(ctx, firstEmail, password)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	userFirst, err := s.domain.user.GetByEmail(ctx, firstEmail)
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}

	res, err := s.domain.session.ListForUser(ctx, userFirst.ID, 0, 10)
	if err != nil {
		t.Fatalf("ListOwnSessions: %v", err)
	}
	if res.Total != 1 || len(res.Data) != 1 {
		t.Fatalf("ListOwnSessions: expected 1 session, got %d", res.Total)
	}

	err = s.domain.session.DeleteAllForUser(ctx, userFirst.ID)
	if err != nil {
		t.Fatalf("DeleteOwnSessions: %v", err)
	}

	res, err = s.domain.session.ListForUser(ctx, userFirst.ID, 0, 10)
	if err != nil {
		t.Fatalf("ListOwnSessions after delete: %v", err)
	}
	if res.Total != 0 || len(res.Data) != 0 {
		t.Fatalf("ListOwnSessions after delete: expected 0 sessions, got %d", res.Total)
	}

	_, err = s.domain.session.Login(ctx, firstEmail, password)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	_, err = s.domain.session.Login(ctx, firstEmail, password)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	res, err = s.domain.session.ListForUser(ctx, userFirst.ID, 0, 10)
	if err != nil {
		t.Fatalf("ListOwnSessions: %v", err)
	}
	if res.Total != 2 || len(res.Data) != 2 {
		t.Fatalf("ListOwnSessions: expected 2 session, got %d", res.Total)
	}
}

func TestUpdateUserPassword(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	firstEmail := "tests@example"
	password := "Test@1234"

	_, err = s.domain.user.Register(ctx,
		firstEmail,
		password,
		roles.User,
	)
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	_, err = s.domain.session.Login(ctx, firstEmail, password)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	userFirst, err := s.domain.user.GetByEmail(ctx, firstEmail)
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}

	res, err := s.domain.session.ListForUser(ctx, userFirst.ID, 0, 10)
	if err != nil {
		t.Fatalf("ListOwnSessions: %v", err)
	}
	if res.Total != 1 || len(res.Data) != 1 {
		t.Fatalf("ListOwnSessions: expected 1 session, got %d", res.Total)
	}

	newPassword := "Test2@1234"

	err = s.domain.user.UpdatePassword(ctx, userFirst.ID, password, newPassword)
	if err != nil {
		t.Fatalf("UpdatePassword: %v", err)
	}

	_, err = s.domain.session.Login(ctx, firstEmail, password)
	if !errors.Is(err, errx.ErrorInvalidLogin) {
		t.Fatalf("Login with old password: expected error %v, got %v", errx.ErrorInvalidLogin, err)
	}

	_, err = s.domain.session.Login(ctx, firstEmail, newPassword)
	if err != nil {
		t.Fatalf("Login with new password: %v", err)
	}
}
