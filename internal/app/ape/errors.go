package ape

import (
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrUserDoesNotExist           = &BusinessError{reason: ReasonUserDoesNotExist}
	ErrSessionDoesNotExist        = &BusinessError{reason: ReasonSessionDoesNotExist}
	ErrUserAlreadyExists          = &BusinessError{reason: ReasonUserAlreadyExists}
	ErrSessionsForUserNotExist    = &BusinessError{reason: ReasonSessionsForUserNotExist}
	ErrSessionClientMismatch      = &BusinessError{reason: ReasonSessionClientMismatch}
	ErrSessionTokenMismatch       = &BusinessError{reason: ReasonSessionTokenMismatch}
	ErrSessionDoesNotBelongToUser = &BusinessError{reason: ReasonSessionDoesNotBelongToUser}
	ErrNoPermission               = &BusinessError{reason: ReasonNoPermission}
	ErrUserSuspended              = &BusinessError{reason: ReasonUserSuspended}
	ErrInternal                   = &BusinessError{reason: ReasonInternal}
)

func ErrorUserDoesNotExist(userID uuid.UUID, err error) error {
	return &BusinessError{
		reason:  ErrUserDoesNotExist.reason,
		message: fmt.Sprintf("user does not exist with ID: %s", userID),
		cause:   err,
	}
}

func ErrorUserDoesNotExistByEmail(email string, cause error) error {
	return &BusinessError{
		reason:  ErrUserDoesNotExist.reason,
		message: fmt.Sprintf("user does not exist with email: %s", email),
		cause:   cause,
	}
}

func ErrorSessionDoesNotExist(sessionID uuid.UUID, cause error) error {
	return &BusinessError{
		reason:  ErrSessionDoesNotExist.reason,
		message: fmt.Sprintf("session does not exist with ID: %s", sessionID),
		cause:   cause,
	}
}

func ErrorUserAlreadyExists(cause error) error {
	return &BusinessError{
		reason:  ErrUserAlreadyExists.reason,
		message: fmt.Sprintf("user already exists"),
		cause:   cause,
	}
}

func ErrorSessionsForUserNotExist(cause error) error {
	return &BusinessError{
		reason:  ErrSessionsForUserNotExist.reason,
		message: fmt.Sprintf("sessions for user do not exist"),
		cause:   cause,
	}
}

func ErrorSessionClientMismatch(cause error) error {
	return &BusinessError{
		reason:  ErrSessionClientMismatch.reason,
		message: "session client mismatch",
		cause:   cause,
	}
}

func ErrorSessionTokenMismatch(cause error) error {
	return &BusinessError{
		reason:  ErrSessionTokenMismatch.reason,
		message: "session token mismatch",
		cause:   cause,
	}
}

func ErrorSessionDoesNotBelongToUser(sessionID, userID uuid.UUID) error {
	return &BusinessError{
		reason:  ErrSessionDoesNotBelongToUser.reason,
		message: fmt.Sprintf("session %s does not belong to user %s", sessionID, userID),
		cause:   fmt.Errorf("session %s does not belong to user %s", sessionID, userID),
	}
}

func ErrorNoPermission(cause error) error {
	return &BusinessError{
		reason:  ErrNoPermission.reason,
		message: fmt.Sprintf("error no permision for this"),
		cause:   cause,
	}
}

func ErrorUserSuspended(userID uuid.UUID) error {
	return &BusinessError{
		reason:  ErrUserSuspended.reason,
		message: fmt.Sprintf("user %s is suspended", userID),
		cause:   fmt.Errorf("user %s is suspended", userID),
	}
}

func ErrorInternal(cause error) error {
	return &BusinessError{
		reason:  ErrInternal.reason,
		message: "unexpected internal error occurred",
		cause:   cause,
	}
}
