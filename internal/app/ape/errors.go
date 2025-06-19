package ape

import (
	"fmt"

	"github.com/google/uuid"
)

type Error struct {
	Err   error
	cause error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Err.Error(), e.cause.Error())
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Nil() bool {
	if e == nil {
		return true
	}
	return e.Err == nil && e.cause == nil
}

var ErrUserDoesNotExist = fmt.Errorf("user does not exist")
var ErrSessionDoesNotExist = fmt.Errorf("session does not exist")
var ErrUserAlreadyExists = fmt.Errorf("user already exists")

var ErrSessionsForUserNotExist = fmt.Errorf("sessions for user does not exist")
var ErrSessionClientMismatch = fmt.Errorf("session client mismatch")
var ErrSessionTokenMismatch = fmt.Errorf("session token mismatch")
var ErrSessionDoesNotBelongToUser = fmt.Errorf("session does not belong to user")

var ErrNoPermission = fmt.Errorf("no permission to perform this action")

var ErrInternal = fmt.Errorf("internal server error")

func ErrorUserDoesNotExist(userID uuid.UUID, err error) error {
	return &Error{Err: ErrUserDoesNotExist, cause: err}
}

func ErrorUserDoesNotExistByEmail(email string, cause error) error {
	return &Error{Err: ErrUserDoesNotExist, cause: fmt.Errorf("%s: %w", email, cause)}
}

func ErrorSessionDoesNotExist(sessionID uuid.UUID, cause error) error {
	return &Error{Err: ErrSessionDoesNotExist, cause: fmt.Errorf("%s: %w", sessionID, cause)}
}

func ErrorUserAlreadyExists(cause error) error {
	return &Error{Err: ErrUserAlreadyExists, cause: cause}
}

func ErrorSessionsForUserNotExist(cause error) error {
	return &Error{Err: ErrSessionsForUserNotExist, cause: cause}
}

func ErrorSessionClientMismatch(cause error) error {
	return &Error{Err: ErrSessionClientMismatch, cause: cause}
}

func ErrorSessionTokenMismatch(cause error) error {
	return &Error{Err: ErrSessionTokenMismatch, cause: cause}
}

func ErrorSessionDoesNotBelongToUser(sessionID, userID uuid.UUID) error {
	return &Error{
		Err:   ErrSessionDoesNotBelongToUser,
		cause: fmt.Errorf("session %s does not belong to user %s", sessionID, userID),
	}
}

func ErrorNoPermission(cause error) error {
	return &Error{Err: ErrNoPermission, cause: cause}
}

func ErrorInternal(cause error) error {
	return &Error{Err: ErrInternal, cause: cause}
}
