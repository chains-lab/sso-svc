package password

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/chains-lab/sso-svc/internal/errx"
	"golang.org/x/crypto/bcrypt"
)

func ReliabilityCheck(password string) error {
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

var InvalidCredentialsError = fmt.Errorf("invalid credentials")

func CheckPasswordMatch(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errx.ErrorInvalidLogin.Raise(
				fmt.Errorf("invalid credentials, cause: %w", err),
			)
		}

		return errx.ErrorInternal.Raise(
			fmt.Errorf("comparing password hash, cause: %w", err),
		)
	}

	return nil
}
