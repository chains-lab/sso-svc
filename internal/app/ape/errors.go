package ape

import (
	"fmt"

	"github.com/google/uuid"
)

type Error struct {
	Reason  string // similar to CODE in HTTP API errors
	Details error  // for internal use in application
	cause   error  // the original error that caused this error, if any
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Reason, e.cause.Error())
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Nil() bool {
	if e == nil {
		return true
	}
	return e.Details == nil && e.cause == nil
}

var ErrUserDoesNotExist = fmt.Errorf("user does not exist")
var ErrSessionDoesNotExist = fmt.Errorf("session does not exist")
var ErrUserAlreadyExists = fmt.Errorf("user already exists")
var ErrSessionsForUserNotExist = fmt.Errorf("sessions for user does not exist")
var ErrSessionClientMismatch = fmt.Errorf("session client mismatch")
var ErrSessionTokenMismatch = fmt.Errorf("session token mismatch")
var ErrSessionDoesNotBelongToUser = fmt.Errorf("session does not belong to user")
var ErrNoPermission = fmt.Errorf("no permission to perform this action")
var ErrOnlyUserCanHaveSubscription = fmt.Errorf("only ordinary user can have subscription")
var ErrOnlyOrdinaryUserCanBeVerified = fmt.Errorf("only ordinary user can be verified")
var ErrorOnlyOrdinaryUserCanBeSuspended = fmt.Errorf("only ordinary user can be suspended")
var ErrUserIsSuspended = fmt.Errorf("user is suspended")
var ErrInternal = fmt.Errorf("internal server error")

const ReasonUserDoesNotExist = "USER_DOES_NOT_EXIST"

func ErrorUserDoesNotExist(userID uuid.UUID, err error) error {
	return &Error{
		Reason:  ReasonUserDoesNotExist,
		Details: ErrUserDoesNotExist,
		cause:   err,
	}
}

func ErrorUserDoesNotExistByEmail(email string, cause error) error {
	return &Error{
		Reason:  ReasonUserDoesNotExist,
		Details: ErrUserDoesNotExist,
		cause:   fmt.Errorf("user dosent exist with email: %s: %w", email, cause),
	}
}

const ReasonSessionDoesNotExist = "SESSION_DOES_NOT_EXIST"

func ErrorSessionDoesNotExist(sessionID uuid.UUID, cause error) error {
	return &Error{
		Reason:  ReasonSessionDoesNotExist,
		Details: ErrSessionDoesNotExist,
		cause:   fmt.Errorf("sessions dosen exist%s: %w", sessionID, cause),
	}
}

const ReasonUserAlreadyExists = "USER_ALREADY_EXISTS"

func ErrorUserAlreadyExists(cause error) error {
	return &Error{
		Reason:  ReasonUserAlreadyExists,
		Details: ErrUserAlreadyExists,
		cause:   cause,
	}
}

const ReasonSessionsForUserNotExist = "SESSIONS_FOR_USER_NOT_EXIST"

func ErrorSessionsForUserNotExist(cause error) error {
	return &Error{
		Reason:  ReasonSessionsForUserNotExist,
		Details: ErrSessionsForUserNotExist,
		cause:   cause,
	}
}

const ReasonSessionClientMismatch = "SESSION_CLIENT_MISMATCH"

func ErrorSessionClientMismatch(cause error) error {
	return &Error{
		Reason:  ReasonSessionClientMismatch,
		Details: ErrSessionClientMismatch,
		cause:   cause,
	}
}

const ReasonSessionTokenMismatch = "SESSION_TOKEN_MISMATCH"

func ErrorSessionTokenMismatch(cause error) error {
	return &Error{
		Reason:  ReasonSessionTokenMismatch,
		Details: ErrSessionTokenMismatch,
		cause:   cause,
	}
}

const ReasonSessionDoesNotBelongToUser = "SESSION_DOES_NOT_BELONG_TO_USER"

func ErrorSessionDoesNotBelongToUser(sessionID, userID uuid.UUID) error {
	return &Error{
		Reason:  ReasonSessionDoesNotBelongToUser,
		Details: ErrSessionDoesNotBelongToUser,
		cause:   fmt.Errorf("session %s does not belong to user %s", sessionID, userID),
	}
}

const ReasonNoPermission = "NO_PERMISSION"

func ErrorNoPermission(cause error) error {
	return &Error{
		Reason:  ReasonNoPermission,
		Details: ErrNoPermission,
		cause:   cause,
	}
}

const ReasonOnlyOrdinaryUserCanBeSuspended = "ONLY_ORDINARY_USER_CAN_BE_SUSPENDED"

func ErrorOnlyUserCanHaveSubscription(cause error) error {
	return &Error{
		Reason:  ReasonOnlyOrdinaryUserCanBeSuspended,
		Details: ErrOnlyUserCanHaveSubscription,
		cause:   cause,
	}
}

func ErrorOnlyOrdinaryUserCanBeVerified(cause error) error {
	return &Error{
		Reason:  ReasonOnlyOrdinaryUserCanBeSuspended,
		Details: ErrOnlyOrdinaryUserCanBeVerified,
		cause:   cause,
	}
}

const ReasonUserIsSuspended = "USER_IS_SUSPENDED"

func ErrorUserSuspended(userID uuid.UUID) error {
	return &Error{
		Reason:  ReasonUserIsSuspended,
		Details: ErrUserIsSuspended,
		cause:   fmt.Errorf("user %s is suspended", userID),
	}
}

const ReasonInternal = "INTERNAL_ERROR"

func ErrorInternal(cause error) error {
	return &Error{
		Reason:  ReasonInternal,
		Details: ErrInternal,
		cause:   cause,
	}
}
