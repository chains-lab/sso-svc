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

func ErrorAccountDoesNotExistByID(accountID uuid.UUID, err error) *Error {
	return &Error{
		Code:    CodeAccountDoesNotExist,
		Title:   "Account does not exist",
		Details: fmt.Sprintf("The requested account %s does not exist.", accountID),
		cause:   err,
	}
}

func ErrorAccountDoesNotExistByEmail(email string, err error) *Error {
	return &Error{
		Code:    CodeAccountDoesNotExist,
		Title:   "Account does not exist",
		Details: fmt.Sprintf("The requested account %s does not exist.", email),
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

func ErrorAccountAlreadyExists(cause error) *Error {
	return &Error{
		Code:    CodeAccountAlreadyExists,
		Title:   "Account Already Exists",
		Details: "Account already exists",
		cause:   cause,
	}
}

func ErrorAccountInvalidRole(cause error) *Error {
	return &Error{
		Code:    CodeAccountInvalidRole,
		Title:   "Invalid Account Role",
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

func ErrorSessionsForAccountNotExist(cause error) *Error {
	return &Error{
		Code:    CodeSessionsForAccountNotExist,
		Title:   "Sessions for Account Not Found",
		Details: "Sessions for account does not exist",
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

func ErrorSessionCannotBeCurrentAccount(cause error) *Error {
	return &Error{
		Code:    CodeSessionCannotBeCurrentAccount,
		Title:   "Session Cannot Be Current Account",
		Details: "Session cannot be current account",
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

func ErrorInternalServer(cause error) *Error {
	return &Error{
		Code:    CodeInternal,
		Title:   "Internal Server Error",
		Details: "Internal server error",
		cause:   cause,
	}
}
