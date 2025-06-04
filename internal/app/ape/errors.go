package ape

import (
	"fmt"

	"github.com/google/uuid"
)

type Error struct {
	Code    string
	Title   string
	Details string
	cause   error
}

func (e *Error) Unwrap() error {
	return e.cause
}

func ErrorUserDoesNotExist(userID uuid.UUID, err error) *Error {
	return &Error{
		Code:    CodeUserDoesNotExist,
		Title:   "User does not exist",
		Details: fmt.Sprintf("The requested user %s does not exist.", userID),
		cause:   err,
	}
}

func ErrorUserDoesNotExistByEmail(email string, err error) *Error {
	return &Error{
		Code:    CodeUserDoesNotExist,
		Title:   "User does not exist",
		Details: fmt.Sprintf("The requested user with email %s does not exist.", email),
		cause:   err,
	}
}

func ErrorSessionDoesNotExist(sessionID uuid.UUID, err error) *Error {
	return &Error{
		Code:    CodeSessionDoesNotExist,
		Title:   "Session does not exist",
		Details: fmt.Sprintf("The requested session %s does not exist.", sessionID),
		cause:   err,
	}
}

func ErrorUserAlreadyExists(cause error) *Error {
	return &Error{
		Code:    CodeUserAlreadyExists,
		Title:   "User Already Exists",
		Details: "User already exists",
		cause:   cause,
	}
}

func ErrorUserInvalidRole(cause error) *Error {
	return &Error{
		Code:    CodeUserInvalidRole,
		Title:   "Invalid User Role",
		Details: "Invalid role",
		cause:   cause,
	}
}

func ErrorUserNoPermissionToUpdateRole(cause error) *Error {
	return &Error{
		Code:    CodeUserNoPermissionToUpdateRole,
		Title:   "No Permission to Update Role",
		Details: "User has no permission to update role",
		cause:   cause,
	}
}

func ErrorSessionsForUserNotExist(cause error) *Error {
	return &Error{
		Code:    CodeSessionsForUserNotExist,
		Title:   "Sessions for User Not Found",
		Details: "Sessions for user does not exist",
		cause:   cause,
	}
}

func ErrorSessionClientMismatch(cause error) *Error {
	return &Error{
		Code:    CodeSessionClientMismatch,
		Title:   "Sessions Client Mismatch",
		Details: "Client does not match",
		cause:   cause,
	}
}

func ErrorSessionTokenMismatch(cause error) *Error {
	return &Error{
		Code:    CodeSessionTokenMismatch,
		Title:   "Sessions Token Mismatch",
		Details: "Token does not match",
		cause:   cause,
	}
}

func ErrorSessionAlreadyExists(cause error) *Error {
	return &Error{
		Code:    CodeSessionAlreadyExists,
		Title:   "Session Already Exists",
		Details: "Session already exists",
		cause:   cause,
	}
}

func ErrorSessionCannotBeCurrent(cause error) *Error {
	return &Error{
		Code:    CodeSessionCannotBeCurrent,
		Title:   "Session Cannot Be Current",
		Details: "Session cannot be current",
		cause:   cause,
	}
}

func ErrorSessionCannotBeCurrentUser(cause error) *Error {
	return &Error{
		Code:    CodeSessionCannotBeCurrentUser,
		Title:   "Session Cannot Be Current User",
		Details: "Session cannot be current user",
		cause:   cause,
	}
}

func ErrorSessionCannotDeleteSuperUserByOther(cause error) *Error {
	return &Error{
		Code:    CodeSessionCannotDeleteSuperUserByOther,
		Title:   "Cannot Delete Superuser Session",
		Details: "Cannot delete superuser session by other user",
		cause:   cause,
	}
}

func ErrorInternal(cause error) *Error {
	return &Error{
		Code:    CodeInternal,
		Title:   "Internal Server Error",
		Details: "Internal server error",
		cause:   cause,
	}
}
