package errx

import (
	"fmt"

	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorUserNotFound = ape.Declare("USER_NOT_FOUND")

func RaiseUserNotFound(cause error, userID uuid.UUID) error {
	return ErrorUserNotFound.Raise(
		cause,
		status.New(codes.NotFound, fmt.Sprintf("user with id: %s not found", userID)),
	)
}

func RaiseUserNotFoundByEmail(cause error, email string) error {
	return ErrorUserNotFound.Raise(
		cause,
		status.New(codes.NotFound, fmt.Sprintf("user with email: %s not found", email)),
	)
}

var ErrorUserAlreadyExists = ape.Declare("USER_ALREADY_EXISTS")

func RaiseUserAlreadyExists(cause error, email string) error {
	return ErrorUserAlreadyExists.Raise(
		cause,
		status.New(codes.AlreadyExists, fmt.Sprintf("user with email: %s already exists", email)),
	)
}

var ErrorUserSuspended = ape.Declare("USER_SUSPENDED")

func RaiseUserSuspended(cause error, userID uuid.UUID) error {
	return ErrorUserSuspended.Raise(
		cause,
		status.New(codes.FailedPrecondition, fmt.Sprintf("user with id: %s is suspended", userID)),
	)
}

var ErrorInitiatorUserSuspended = ape.Declare("USER_INITIATOR_SUSPENDED")

func RaiseInitiatorUserSuspended(cause error, userID uuid.UUID) error {
	return ErrorInitiatorUserSuspended.Raise(
		cause,
		status.New(codes.PermissionDenied, fmt.Sprintf("initiator with id: %s is suspended", userID)),
	)
}

var ErrorUserRoleIsNotAllowed = ape.Declare("USER_ROLE_IS_NOT_ALLOWED")

func RaiseUserRoleIsNotAllowed(cause error) error {
	return ErrorUserRoleIsNotAllowed.Raise(
		cause,
		status.New(codes.PermissionDenied, cause.Error()),
	)
}

var ErrorInitiatorRoleIsLowThanTarget = ape.Declare("INITIATOR_ROLE_IS_LOWER_THAN_TARGET")

func RaiseInitiatorRoleIsLowThanTarget(cause error) error {
	return ErrorInitiatorRoleIsLowThanTarget.Raise(
		cause,
		status.New(codes.PermissionDenied, cause.Error()),
	)
}

var ErrorInitiatorNotFound = ape.Declare("INITIATOR_NOT_FOUND")

func RaiseInitiatorNotFound(cause error, initiatorID uuid.UUID) error {
	return ErrorInitiatorNotFound.Raise(
		cause,
		status.New(codes.NotFound, fmt.Sprintf("initiator with id: %s not found", initiatorID)),
	)
}
