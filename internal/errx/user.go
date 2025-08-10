package errx

import (
	"fmt"

	"github.com/chains-lab/sso-svc/internal/errx/statusx"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
)

var ErrorUserNotFound = ape.Declare("USER_NOT_FOUND")

func RaiseUserNotFound(cause error, userID uuid.UUID) error {
	return ErrorUserNotFound.Raise(
		cause,
		statusx.NotFound(fmt.Sprintf("user with id: %s not found", userID)),
	)
}

func RaiseUserNotFoundByEmail(cause error, email string) error {
	return ErrorUserNotFound.Raise(
		cause,
		statusx.NotFound(fmt.Sprintf("user with email: %s not found", email)),
	)
}

var ErrorUserAlreadyExists = ape.Declare("USER_ALREADY_EXISTS")

func RaiseUserAlreadyExists(cause error, email string) error {
	return ErrorUserAlreadyExists.Raise(
		cause,
		statusx.AlreadyExists(fmt.Sprintf("user with email: %s already exists", email)),
	)
}

var ErrorUserSuspended = ape.Declare("USER_SUSPENDED")

func RaiseUserSuspended(cause error, userID uuid.UUID) error {
	return ErrorUserSuspended.Raise(
		cause,
		statusx.PermissionDenied(fmt.Sprintf("user with id: %s is suspended", userID)),
	)
}
