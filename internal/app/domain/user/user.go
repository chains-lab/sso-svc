package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"unicode"

	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/google/uuid"
)

type usersQ interface {
	New() dbx.UserQ
	Insert(ctx context.Context, input dbx.UserModel) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]dbx.UserModel, error)
	Get(ctx context.Context) (dbx.UserModel, error)

	FilterID(id uuid.UUID) dbx.UserQ
	FilterRole(role string) dbx.UserQ

	Update(ctx context.Context, input map[string]any) error

	Page(limit, offset uint64) dbx.UserQ
	Count(ctx context.Context) (uint64, error)
}

type emailQ interface {
	New() dbx.UserEmailQ
	Insert(ctx context.Context, input dbx.UserEmailModel) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]dbx.UserEmailModel, error)
	Get(ctx context.Context) (dbx.UserEmailModel, error)

	FilterID(id uuid.UUID) dbx.UserEmailQ
	FilterEmail(email string) dbx.UserEmailQ

	Update(ctx context.Context, input map[string]any) error

	Page(limit, offset uint64) dbx.UserEmailQ
	Count(ctx context.Context) (uint64, error)
}

type passQ interface {
	New() dbx.UserPassQ
	Insert(ctx context.Context, input dbx.UserPasswordModel) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]dbx.UserPasswordModel, error)
	Get(ctx context.Context) (dbx.UserPasswordModel, error)

	FilterID(id uuid.UUID) dbx.UserPassQ

	Update(ctx context.Context, input map[string]any) error

	Page(limit, offset uint64) dbx.UserPassQ
	Count(ctx context.Context) (uint64, error)
}

type User struct {
	query  usersQ
	passQ  passQ
	emailQ emailQ
}

func CreateUser(pg *sql.DB) User {
	return User{
		query:  dbx.NewUsers(pg),
		passQ:  dbx.NewUsersPass(pg),
		emailQ: dbx.NewUsersEmail(pg),
	}
}

func checkPassword(password string) error {
	if len(password) < 8 || len(password) > 32 {
		return fmt.Errorf("password must be between 8 and 32 characters")
	}

	var (
		hasUpper, hasLower, hasDigit, hasSpecial bool
	)

	allowedSpecials := "-.!#$%&?,@"

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case strings.ContainsRune(allowedSpecials, r):
			hasSpecial = true
		default:
			return fmt.Errorf("password contains invalid characters %s", string(r))
		}
	}

	if !hasUpper {
		return fmt.Errorf("need at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("need at least one lower case letter")
	}
	if !hasDigit {
		return fmt.Errorf("need at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("need at least one special character from %s", allowedSpecials)
	}

	return nil
}
