package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/errx"
)

func TestAdminBlockUser(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	superAdmin := CreateUser(s, t, "superadmin@example", "Super@1234", roles.SuperUser)
	superAdminSes := CreateSession(s, t, superAdmin.ID)
	admin := CreateUser(s, t, "admin@example", "Admin@1234", roles.Admin)
	adminSes := CreateSession(s, t, admin.ID)
	_ = CreateSession(s, t, admin.ID)
	_ = CreateSession(s, t, admin.ID)
	user := CreateUser(s, t, "user@example", "User@1234", roles.User)
	userSes := CreateSession(s, t, user.ID)

	_, err = s.app.AdminBlockUser(ctx, admin.ID, adminSes.ID, superAdmin.ID)
	if !errors.Is(err, errx.ErrorNoPermissions) {
		t.Fatalf("AdminBlockUser: expected no permissions error, got %v", err)
	}

	_, err = s.app.AdminBlockUser(ctx, user.ID, userSes.ID, superAdmin.ID)
	if !errors.Is(err, errx.ErrorNoPermissions) {
		t.Fatalf("AdminBlockUser: expected no permissions error, got %v", err)
	}

	_, err = s.app.AdminBlockUser(ctx, superAdmin.ID, superAdminSes.ID, superAdmin.ID)
	if !errors.Is(err, errx.ErrorUserCannotBlockHimself) {
		t.Fatalf("AdminBlockUser: expected no permissions error, got %v", err)
	}

	sess, _, err := s.app.ListOwnSessions(ctx, admin.ID, pagi.Request{}, nil)
	if err != nil {
		t.Fatalf("ListOwnSessions: unexpected error: %v", err)
	}
	if len(sess) != 3 {
		t.Fatalf("ListOwnSessions: expected 3 sessions, got %d", len(sess))
	}

	admin, err = s.app.AdminBlockUser(ctx, superAdmin.ID, superAdminSes.ID, admin.ID)
	if err != nil {
		t.Fatalf("AdminBlockUser: unexpected error: %v", err)
	}

	if admin.Status != enum.UserStatusBlocked {
		t.Fatalf("AdminBlockUser: expected status 'blocked', got %s", admin.Status)
	}

	admin, err = s.app.AdminUnblockUser(ctx, superAdmin.ID, superAdminSes.ID, admin.ID)
	if err != nil {
		t.Fatalf("AdminUnblockUser: unexpected error: %v", err)
	}

	if admin.Status != enum.UserStatusActive {
		t.Fatalf("AdminUnblockUser: expected status 'active', got %s", admin.Status)
	}

	sess, _, err = s.app.ListOwnSessions(ctx, admin.ID, pagi.Request{}, nil)
	if err != nil {
		t.Fatalf("ListOwnSessions: unexpected error: %v", err)
	}
	if len(sess) != 0 {
		t.Fatalf("ListOwnSessions: expected 0 sessions, got %d", len(sess))
	}

	_ = CreateSession(s, t, admin.ID)

	sess, _, err = s.app.ListOwnSessions(ctx, admin.ID, pagi.Request{}, nil)
	if err != nil {
		t.Fatalf("ListOwnSessions: unexpected error: %v", err)
	}
	if len(sess) != 1 {
		t.Fatalf("ListOwnSessions: expected 1 sessions, got %d", len(sess))
	}

	_ = CreateSession(s, t, user.ID)
	_ = CreateSession(s, t, user.ID)
	_ = CreateSession(s, t, user.ID)
	sesUserFive := CreateSession(s, t, user.ID)

	sess, _, err = s.app.ListOwnSessions(ctx, user.ID, pagi.Request{}, nil)
	if err != nil {
		t.Fatalf("ListOwnSessions: unexpected error: %v", err)
	}
	if len(sess) != 5 {
		t.Fatalf("ListOwnSessions: expected 5 sessions, got %d", len(sess))
	}

	sess, _, err = s.app.AdminListUserSessions(ctx, superAdmin.ID, superAdminSes.ID, user.ID, pagi.Request{}, nil)
	if err != nil {
		t.Fatalf("AdminListUserSessions: unexpected error: %v", err)
	}
	if len(sess) != 5 {
		t.Fatalf("AdminListUserSessions: expected 5 sessions, got %d", len(sess))
	}

	err = s.app.AdminDeleteUserSession(ctx, superAdmin.ID, superAdminSes.ID, user.ID, sesUserFive.ID)
	if err != nil {
		t.Fatalf("AdminDeleteUserSession: unexpected error: %v", err)
	}

	_, err = s.app.GetOwnSession(ctx, user.ID, sesUserFive.ID)
	if !errors.Is(err, errx.ErrorSessionNotFound) {
		t.Fatalf("GetOwnSession: expected session not found error, got %v", err)
	}

	sess, _, err = s.app.ListOwnSessions(ctx, user.ID, pagi.Request{}, nil)
	if err != nil {
		t.Fatalf("ListOwnSessions: unexpected error: %v", err)
	}
	if len(sess) != 4 {
		t.Fatalf("ListOwnSessions: expected 4 sessions, got %d", len(sess))
	}
	for _, s := range sess {
		if s.ID == sesUserFive.ID {
			t.Fatalf("AdminDeleteUserSession: session %s should be deleted", sesUserFive.ID)
		}
	}

	err = s.app.AdminDeleteUserSessions(ctx, superAdmin.ID, superAdminSes.ID, user.ID)
	if err != nil {
		t.Fatalf("AdminDeleteUserSessions: unexpected error: %v", err)
	}

	sess, _, err = s.app.ListOwnSessions(ctx, user.ID, pagi.Request{}, nil)
	if err != nil {
		t.Fatalf("ListOwnSessions: unexpected error: %v", err)
	}
	if len(sess) != 0 {
		t.Fatalf("ListOwnSessions: expected 0 sessions, got %d", len(sess))
	}
}
