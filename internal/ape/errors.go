package ape

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ServiceName = "sso-svc"

var (
	ErrUserDoesNotExist           = &Error{reason: ReasonUserDoesNotExist}
	ErrSessionDoesNotExist        = &Error{reason: ReasonSessionDoesNotExist}
	ErrUserAlreadyExists          = &Error{reason: ReasonUserAlreadyExists}
	ErrSessionsForUserNotExist    = &Error{reason: ReasonSessionsForUserNotExist}
	ErrSessionClientMismatch      = &Error{reason: ReasonSessionClientMismatch}
	ErrSessionTokenMismatch       = &Error{reason: ReasonSessionTokenMismatch}
	ErrSessionDoesNotBelongToUser = &Error{reason: ReasonSessionDoesNotBelongToUser}
	ErrNoPermission               = &Error{reason: ReasonNoPermission}
	ErrUserSuspended              = &Error{reason: ReasonUserSuspended}
	ErrInternal                   = &Error{reason: ReasonInternal}
)

func RaiseUserNotFound(userID uuid.UUID, cause error) error {
	msg := fmt.Sprintf("user %s does not exist", userID)
	return &Error{
		code:    codes.NotFound,
		reason:  ErrUserDoesNotExist.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{Reason: ErrUserDoesNotExist.reason, Domain: ServiceName},
			&errdetails.ResourceInfo{ResourceType: "user", ResourceName: userID.String(), Description: msg},
		},
	}
}

// RaiseUserNotFoundByEmail возвращает NotFound для отсутствующего пользователя по email
func RaiseUserNotFoundByEmail(email string, cause error) error {
	msg := fmt.Sprintf("user with email %s does not exist", email)
	return &Error{
		code:    codes.NotFound,
		reason:  ErrUserDoesNotExist.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{Reason: ErrUserDoesNotExist.reason, Domain: ServiceName},
			&errdetails.ResourceInfo{ResourceType: "user_email", ResourceName: email, Description: msg},
		},
	}
}

// RaiseUserAlreadyExists возвращает AlreadyExists, когда пользователь уже существует
func RaiseUserAlreadyExists(cause error) error {
	msg := "user already exists"
	return &Error{
		code:    codes.AlreadyExists,
		reason:  ErrUserAlreadyExists.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{Reason: ErrUserAlreadyExists.reason, Domain: ServiceName},
		},
	}
}

// RaiseSessionNotFound возвращает NotFound для отсутствующей сессии
func RaiseSessionNotFound(sessionID uuid.UUID, cause error) error {
	msg := fmt.Sprintf("session %s does not exist", sessionID)
	return &Error{
		code:    codes.NotFound,
		reason:  ErrSessionDoesNotExist.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{Reason: ErrSessionDoesNotExist.reason, Domain: ServiceName},
			&errdetails.ResourceInfo{ResourceType: "session", ResourceName: sessionID.String(), Description: msg},
		},
	}
}

// RaiseSessionsForUserNotExist возвращает NotFound, если у пользователя нет сессий
func RaiseSessionsForUserNotExist(userID uuid.UUID, cause error) error {
	msg := fmt.Sprintf("no sessions found for user %s", userID)
	return &Error{
		code:    codes.NotFound,
		reason:  ErrSessionsForUserNotExist.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{Reason: ErrSessionsForUserNotExist.reason, Domain: ServiceName},
			&errdetails.ResourceInfo{ResourceType: "session_list", ResourceName: userID.String(), Description: msg},
		},
	}
}

// RaiseSessionClientMismatch возвращает FailedPrecondition при несоответствии клиента сессии
func RaiseSessionClientMismatch(cause error) error {
	msg := "session client mismatch"
	return &Error{
		code:    codes.FailedPrecondition,
		reason:  ErrSessionClientMismatch.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{Reason: ErrSessionClientMismatch.reason, Domain: ServiceName},
		},
	}
}

// RaiseSessionTokenMismatch возвращает FailedPrecondition при несоответствии токена
func RaiseSessionTokenMismatch(cause error) error {
	msg := "session token mismatch"
	return &Error{
		code:    codes.FailedPrecondition,
		reason:  ErrSessionTokenMismatch.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{Reason: ErrSessionTokenMismatch.reason, Domain: ServiceName},
		},
	}
}

// RaiseSessionDoesNotBelongToUser возвращает PermissionDenied, если сессия чужая
func RaiseSessionDoesNotBelongToUser(sessionID, userID uuid.UUID) error {
	msg := fmt.Sprintf("session %s does not belong to user %s", sessionID, userID)
	return &Error{
		code:    codes.PermissionDenied,
		reason:  ErrSessionDoesNotBelongToUser.reason,
		message: msg,
		cause:   fmt.Errorf("session %s does not belong to user %s", sessionID, userID),
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{Reason: ErrSessionDoesNotBelongToUser.reason, Domain: ServiceName},
			&errdetails.ResourceInfo{ResourceType: "session", ResourceName: sessionID.String(), Description: msg},
		},
	}
}

// RaiseNoPermission возвращает PermissionDenied при отсутствии прав
func RaiseNoPermission(cause error) error {
	msg := "no permission to perform this action"
	return &Error{
		code:    codes.PermissionDenied,
		reason:  ErrNoPermission.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{Reason: ErrNoPermission.reason, Domain: ServiceName},
		},
	}
}

// RaiseUserSuspended возвращает FailedPrecondition при блокировке пользователя
func RaiseUserSuspended(userID uuid.UUID) error {
	msg := fmt.Sprintf("user %s is suspended", userID)
	return &Error{
		code:    codes.FailedPrecondition,
		reason:  ErrUserSuspended.reason,
		message: msg,
		cause:   fmt.Errorf("user %s is suspended", userID),
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{Reason: ErrUserSuspended.reason, Domain: ServiceName},
			&errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{{Type: "user_suspended", Subject: userID.String(), Description: msg}}},
		},
	}
}

// RaiseInternal возвращает базовую Internal ошибку
func RaiseInternal(cause error) error {
	return &Error{
		code:    codes.Internal,
		reason:  ErrInternal.reason,
		message: "unexpected internal error occurred",
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ErrorInfo{Reason: ErrInternal.reason, Domain: ServiceName},
		},
	}
}
