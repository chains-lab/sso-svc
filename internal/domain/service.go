package domain

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/chains-lab/restkit/token"
	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/google/uuid"
)

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	ParseRefreshClaims(enc string) (token.AccountClaims, error)

	GenerateAccess(
		account entity.Account, sessionID uuid.UUID,
	) (string, error)

	GenerateRefresh(
		account entity.Account, sessionID uuid.UUID,
	) (string, error)
}

type EventPublisher interface {
	WriteAccountCreated(ctx context.Context, account entity.Account, email string) error
	WriteAccountPasswordChanged(ctx context.Context, account entity.Account, email string) error
	WriteAccountUsernameChanged(ctx context.Context, account entity.Account, email string) error
	WriteAccountLogin(ctx context.Context, account entity.Account, email string) error
}

type CreateAccountParams struct {
	Username     string
	Role         string
	Email        string
	PasswordHash string
}

type database interface {
	CreateAccount(
		ctx context.Context,
		params CreateAccountParams,
	) (entity.Account, error)

	GetAccountByID(ctx context.Context, accountID uuid.UUID) (entity.Account, error)
	GetAccountByUsername(ctx context.Context, username string) (entity.Account, error)
	UpdateAccountUsername(
		ctx context.Context,
		accountID uuid.UUID,
		newUsername string,
	) (entity.Account, error)

	GetAccountByEmail(ctx context.Context, email string) (entity.Account, error)
	UpdateAccountStatus(
		ctx context.Context,
		accountID uuid.UUID,
		status string,
	) (entity.Account, error)

	GetAccountEmail(ctx context.Context, accountID uuid.UUID) (entity.AccountEmail, error)
	UpdateAccountEmailVerification(
		ctx context.Context,
		accountID uuid.UUID,
		verified bool,
	) (entity.AccountEmail, error)

	GetAccountPassword(ctx context.Context, accountID uuid.UUID) (entity.AccountPassword, error)
	UpdateAccountPassword(
		ctx context.Context,
		accountID uuid.UUID,
		passwordHash string,
	) error
	DeleteAccount(ctx context.Context, accountID uuid.UUID) error

	CreateSession(ctx context.Context, sessionID, accountID uuid.UUID, hashToken string) (entity.Session, error)
	GetSession(ctx context.Context, sessionID uuid.UUID) (entity.Session, error)
	GetAccountSession(
		ctx context.Context,
		accountID, sessionID uuid.UUID,
	) (entity.Session, error)
	GetSessionsForAccount(
		ctx context.Context,
		accountID uuid.UUID,
		page, size int32,
	) (entity.SessionsCollection, error)
	GetSessionToken(ctx context.Context, sessionID uuid.UUID) (string, error)
	UpdateSessionToken(
		ctx context.Context,
		sessionID uuid.UUID,
		token string,
	) (entity.Session, error)

	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	DeleteSessionsForAccount(ctx context.Context, accountID uuid.UUID) error
	DeleteAccountSession(ctx context.Context, accountID, sessionID uuid.UUID) error
}

type Service struct {
	db    database
	jwt   JWTManager
	event EventPublisher
}

func (s Service) CheckPasswordRequirements(password string) error {
	if len(password) < 8 || len(password) > 32 {
		return errx.ErrorPasswordIsNotAllowed.Raise(
			fmt.Errorf("password must be between 8 and 32 characters"),
		)
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
			return errx.ErrorPasswordIsNotAllowed.Raise(
				fmt.Errorf("password contains invalid characters %s", string(r)),
			)
		}
	}

	if !hasUpper {
		return errx.ErrorPasswordIsNotAllowed.Raise(
			fmt.Errorf("need at least one uppercase letter"),
		)
	}
	if !hasLower {
		return errx.ErrorPasswordIsNotAllowed.Raise(
			fmt.Errorf("need at least one lower case letter"),
		)
	}
	if !hasDigit {
		return errx.ErrorPasswordIsNotAllowed.Raise(
			fmt.Errorf("need at least one digit"),
		)
	}
	if !hasSpecial {
		return errx.ErrorPasswordIsNotAllowed.Raise(
			fmt.Errorf("need at least one special character from %s", allowedSpecials),
		)
	}

	return nil
}

func (s Service) CheckUsernameRequirements(username string) error {
	if len(username) < 3 || len(username) > 32 {
		return errx.ErrorUsernameIsNotAllowed.Raise(
			fmt.Errorf("username must be between 3 and 32 characters"),
		)
	}

	for _, r := range username {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-') {
			return errx.ErrorUsernameIsNotAllowed.Raise(
				fmt.Errorf("username contains invalid characters %s", string(r)),
			)
		}
	}

	return nil
}

func NewService(
	db database,
	jwt JWTManager,
	event EventPublisher,
) *Service {
	return &Service{
		db:    db,
		jwt:   jwt,
		event: event,
	}
}
