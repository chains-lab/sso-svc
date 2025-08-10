package errx

import (
	"fmt"

	"github.com/chains-lab/sso-svc/internal/errx/statusx"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
)

var ErrorSessionNotFound = ape.Declare("SESSION_NOT_FOUND")

func RaiseSessionNotFound(cause error, sessionID uuid.UUID) error {
	return ErrorSessionNotFound.Raise(
		cause,
		statusx.NotFound(fmt.Sprintf("session with id: %s not found", sessionID)),
	)
}

var ErrorSessionsForUserNotFound = ape.Declare("SESSIONS_FOR_USER_NOT_FOUND")

func RaiseSessionsForUserNotFound(cause error, userID uuid.UUID) error {
	return ErrorSessionsForUserNotFound.Raise(
		cause,
		statusx.NotFound(fmt.Sprintf("sessions for user with id: %s not found", userID)),
	)
}

var ErrorSessionTokenMismatch = ape.Declare("SESSION_TOKEN_MISMATCH")

func RaiseSessionTokenMismatch(cause error) error {
	return ErrorSessionTokenMismatch.Raise(
		cause,
		statusx.PermissionDenied("session token mismatch"),
	)
}

var ErrorSessionClientMismatch = ape.Declare("SESSION_CLIENT_MISMATCH")

func RaiseSessionClientMismatch(cause error) error {
	return ErrorSessionClientMismatch.Raise(
		cause,
		statusx.FailedPrecondition("session client mismatch"),
	)
}

var ErrorSessionDoesNotBelongToUser = ape.Declare("SESSION_DOES_NOT_BELONG_TO_USER")

func RaiseSessionDoesNotBelongToUser(cause error, sessionID, userID uuid.UUID) error {
	return ErrorSessionDoesNotBelongToUser.Raise(
		cause,
		statusx.FailedPrecondition(fmt.Sprintf("session with id: %s user_id: %s not found", sessionID, userID)),
	)
}
