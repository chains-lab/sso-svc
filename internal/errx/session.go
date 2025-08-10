package errx

import (
	"fmt"

	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorSessionNotFound = ape.Declare("SESSION_NOT_FOUND")

func RaiseSessionNotFound(cause error, sessionID uuid.UUID) error {
	return ErrorSessionNotFound.Raise(
		cause,
		status.New(codes.NotFound, fmt.Sprintf("session with id: %s not found", sessionID)),
	)
}

var ErrorSessionsForUserNotFound = ape.Declare("SESSIONS_FOR_USER_NOT_FOUND")

func RaiseSessionsForUserNotFound(cause error, userID uuid.UUID) error {
	return ErrorSessionsForUserNotFound.Raise(
		cause,
		status.New(codes.NotFound, fmt.Sprintf("sessions for user with id: %s not found", userID)),
	)
}

var ErrorSessionTokenMismatch = ape.Declare("SESSION_TOKEN_MISMATCH")

func RaiseSessionTokenMismatch(cause error) error {
	return ErrorSessionTokenMismatch.Raise(
		cause,
		status.New(codes.PermissionDenied, "session token mismatch"),
	)
}

var ErrorSessionClientMismatch = ape.Declare("SESSION_CLIENT_MISMATCH")

func RaiseSessionClientMismatch(cause error) error {
	return ErrorSessionClientMismatch.Raise(
		cause,
		status.New(codes.FailedPrecondition, "session client mismatch"),
	)
}

var ErrorSessionDoesNotBelongToUser = ape.Declare("SESSION_DOES_NOT_BELONG_TO_USER")

func RaiseSessionDoesNotBelongToUser(cause error, sessionID, userID uuid.UUID) error {
	return ErrorSessionDoesNotBelongToUser.Raise(
		cause,
		status.New(codes.FailedPrecondition, fmt.Sprintf("session with id: %s user_id: %s not found", sessionID, userID)),
	)
}
