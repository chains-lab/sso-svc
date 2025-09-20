package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/errx"
)

func TestUserRegistration(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	firstEmail := "test@example"
	password := "Test@1234"

	err = s.app.Register(ctx, firstEmail, password)
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	_, err = s.app.Login(ctx, firstEmail, password)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	userFirst, err := s.app.GetUserByEmail(ctx, firstEmail)
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}

	res, pag, err := s.app.ListOwnSessions(ctx, userFirst.ID, pagi.Request{Page: 0, Size: 10}, nil)
	if err != nil {
		t.Fatalf("ListOwnSessions: %v", err)
	}
	if pag.Total != 1 || len(res) != 1 {
		t.Fatalf("ListOwnSessions: expected 1 session, got %d", pag.Total)
	}

	err = s.app.DeleteOwnSessions(ctx, userFirst.ID, res[0].ID)
	if err != nil {
		t.Fatalf("DeleteOwnSessions: %v", err)
	}

	res, pag, err = s.app.ListOwnSessions(ctx, userFirst.ID, pagi.Request{Page: 0, Size: 10}, nil)
	if err != nil {
		t.Fatalf("ListOwnSessions after delete: %v", err)
	}
	if pag.Total != 0 || len(res) != 0 {
		t.Fatalf("ListOwnSessions after delete: expected 0 sessions, got %d", pag.Total)
	}

	_, err = s.app.Login(ctx, firstEmail, password)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	_, err = s.app.Login(ctx, firstEmail, password)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	res, pag, err = s.app.ListOwnSessions(ctx, userFirst.ID, pagi.Request{Page: 0, Size: 10}, nil)
	if err != nil {
		t.Fatalf("ListOwnSessions: %v", err)
	}
	if pag.Total != 2 || len(res) != 2 {
		t.Fatalf("ListOwnSessions: expected 2 session, got %d", pag.Total)
	}
}

func TestUpdateUserPassword(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	firstEmail := "test@example"
	password := "Test@1234"

	err = s.app.Register(ctx, firstEmail, password)
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	_, err = s.app.Login(ctx, firstEmail, password)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	userFirst, err := s.app.GetUserByEmail(ctx, firstEmail)
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}

	res, pag, err := s.app.ListOwnSessions(ctx, userFirst.ID, pagi.Request{Page: 0, Size: 10}, nil)
	if err != nil {
		t.Fatalf("ListOwnSessions: %v", err)
	}
	if pag.Total != 1 || len(res) != 1 {
		t.Fatalf("ListOwnSessions: expected 1 session, got %d", pag.Total)
	}

	newPassword := "Test2@1234"

	err = s.app.UpdatePassword(ctx, userFirst.ID, res[0].ID, password, newPassword)
	if err != nil {
		t.Fatalf("UpdatePassword: %v", err)
	}

	_, err = s.app.Login(ctx, firstEmail, password)
	if !errors.Is(err, errx.ErrorInvalidLogin) {
		t.Fatalf("Login with old password: expected error %v, got %v", errx.ErrorInvalidLogin, err)
	}

	_, err = s.app.Login(ctx, firstEmail, newPassword)
	if err != nil {
		t.Fatalf("Login with new password: %v", err)
	}
}
