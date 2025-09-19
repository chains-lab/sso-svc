package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (u User) ComparisonRightsForAdmins(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	dif int,
) (initiator, user models.User, err error) {
	initiator, err = u.GetInitiator(ctx, initiatorID)
	if err != nil {
		return initiator, user, err
	}

	user, err = u.GetByID(ctx, userID)
	if err != nil {
		return initiator, user, err
	}

	if initiatorID == userID {
		return initiator, user, nil
	}

	if user.Role != roles.SuperUser {
		allowed, err := roles.CompareRolesUser(initiator.Role, user.Role)
		if err != nil {
			return initiator, user, errx.ErrorRoleNotSupported.Raise(
				fmt.Errorf("comparing roles between initiator %s and user %s, cause: %w", initiator.Role, user.Role, err),
			)
		}

		if allowed < dif {
			return initiator, user, errx.ErrorNoPermissions.Raise(
				fmt.Errorf("initiator Role %s is not allowed to interact with this user", initiator.Role),
			)
		}
	}

	return initiator, user, nil
}

func (u User) CheckPassword(ctx context.Context, userID uuid.UUID, password string) error {
	secret, err := u.passQ.New().FilterID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorUserNotFound.Raise(
				fmt.Errorf("password for user %s not found, cause: %w", userID, err),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("getting password for user %s, cause: %w", userID, err),
			)
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(secret.PassHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errx.ErrorInvalidLogin.Raise(
				fmt.Errorf("invalid credentials for user %s, cause: %w", userID, err),
			)
		}

		return errx.ErrorInternal.Raise(
			fmt.Errorf("comparing password hash for user %s, cause: %w", userID, err),
		)
	}

	return nil
}

func (u User) UpdatePassword(ctx context.Context, userID uuid.UUID, password string) error {
	err := checkPassword(password)
	if err != nil {
		return errx.ErrorPasswordIsInappropriate.Raise(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("hashing new password for user '%s', cause: %w", userID, err),
		)
	}

	err = u.passQ.New().FilterID(userID).Update(ctx,
		map[string]interface{}{
			"password_hash": string(hash),
			"updated_at":    time.Now().UTC(),
		},
	)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("updating password for user '%s', cause: %w", userID, err),
		)
	}

	return nil
}

func (u User) Create(ctx context.Context, email, password, role string) error {
	err := roles.ParseRole(role)
	if err != nil {
		return errx.ErrorRoleNotSupported.Raise(
			fmt.Errorf("parsing role for new user with email '%s', cause: %w", email, err),
		)
	}

	err = checkPassword(password)
	if err != nil {
		return errx.ErrorPasswordIsInappropriate.Raise(err)
	}

	id := uuid.New()
	now := time.Now().UTC()

	err = u.query.New().Insert(ctx, dbx.UserModel{
		ID:        id,
		Role:      role,
		Status:    enum.UserStatusActive,
		CreatedAt: now,
	})

	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("inserting new user with email '%s', cause: %w", email, err),
		)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("hashing password for user '%s', cause: %w", id, err),
		)
	}

	err = u.passQ.New().Insert(ctx, dbx.UserPasswordModel{
		ID:        id,
		PassHash:  string(hash),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("inserting password for new user with email '%s', cause: %w", email, err),
		)
	}

	err = u.emailQ.New().Insert(ctx, dbx.UserEmailModel{
		ID:       id,
		Email:    email,
		Verified: true,
	})
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("inserting email for new user with email '%s', cause: %w", email, err),
		)
	}

	return nil
}

func (u User) GetInitiator(ctx context.Context, userID uuid.UUID) (models.User, error) {
	user, err := u.query.FilterID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUnauthenticated.Raise(
				fmt.Errorf("user with id '%s' not found, cause: %w", userID, err),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with id '%s', cause: %w", userID, err),
			)
		}
	}

	emailData, err := u.emailQ.New().FilterID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUnauthenticated.Raise(
				fmt.Errorf("email for user with id '%s' not found, cause: %w", userID, err),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get email for user with id '%s', cause: %w", userID, err),
			)
		}
	}

	if user.Status == enum.UserStatusBlocked {
		return models.User{}, errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("user with id '%s' is blocked", userID),
		)
	}

	return models.User{
		ID:        user.ID,
		Email:     emailData.Email,
		Role:      user.Role,
		Status:    user.Status,
		EmailVer:  emailData.Verified,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u User) GetByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	user, err := u.query.FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("user with id '%s' not found, cause: %w", ID, err),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with id '%s', cause: %w", ID, err),
			)
		}
	}

	emailData, err := u.emailQ.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("email for user with id '%s' not found, cause: %w", ID, err),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get email for user with id '%s', cause: %w", ID, err),
			)
		}
	}

	return models.User{
		ID:        user.ID,
		Email:     emailData.Email,
		Role:      user.Role,
		Status:    user.Status,
		EmailVer:  emailData.Verified,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u User) GetByEmail(ctx context.Context, email string) (models.User, error) {
	emailData, err := u.emailQ.New().FilterEmail(email).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("user with email '%s' not found, cause: %w", email, err),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with email '%s', cause: %w", email, err),
			)
		}
	}

	user, err := u.query.FilterID(emailData.ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("user with email '%s' not found, cause: %w", email, err),
			)
		default:
			return models.User{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with email '%s', cause: %w", email, err),
			)
		}
	}

	return models.User{
		ID:        user.ID,
		Email:     emailData.Email,
		Role:      user.Role,
		Status:    user.Status,
		EmailVer:  emailData.Verified,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u User) Delete(ctx context.Context, userID uuid.UUID) error {
	err := u.passQ.New().FilterID(userID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(err)
	}

	err = u.query.New().FilterID(userID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(err)
	}

	return nil
}

func (u User) SetStatus(ctx context.Context, userID uuid.UUID, status string) error {
	err := enum.ParseUserStatus(status)
	if err != nil {
		return errx.ErrorUserStatusNotSupported.Raise(
			fmt.Errorf("parsing status for user %s, cause: %w", userID, err),
		)
	}

	err = u.query.New().FilterID(userID).Update(ctx,
		map[string]interface{}{"status": status},
	)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("updating status for user %s, cause: %w", userID, err),
		)
	}

	return nil
}
