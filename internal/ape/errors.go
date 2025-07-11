package ape

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

const ServiceName = "sso-svc"

var ErrInternal = &Error{reason: ReasonInternal}

func RaiseInternal(cause error) error {
	return &Error{
		code:    codes.Internal,
		reason:  ErrInternal.reason,
		message: "unexpected internal error occurred",
		cause:   cause,
	}
}

var ErrUserNotFound = &Error{reason: ReasonUserNotFound, code: codes.NotFound}

func RaiseUserNotFound(userID uuid.UUID, cause error) error {
	msg := fmt.Sprintf("user %s not found", userID)
	return &Error{
		code:    ErrUserNotFound.code,
		reason:  ErrUserNotFound.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "user",
				ResourceName: fmt.Sprintf("user:id:%s", userID),
				Description:  msg,
			},
		},
	}
}

func RaiseUserNotFoundByEmail(email string, cause error) error {
	msg := fmt.Sprintf("user with email %s does not exist", email)
	return &Error{
		code:    ErrUserNotFound.code,
		reason:  ErrUserNotFound.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "user",
				ResourceName: fmt.Sprintf("user:email:%s", email),
				Description:  msg,
			},
		},
	}
}

var ErrUserAlreadyExists = &Error{reason: ReasonUserAlreadyExists, code: codes.AlreadyExists}

func RaiseUserAlreadyExists(cause error) error {
	msg := "user already exists"
	return &Error{
		code:    ErrUserAlreadyExists.code,
		reason:  ErrUserAlreadyExists.reason,
		message: msg,
		cause:   cause,
	}
}

var ErrSessionNotFound = &Error{reason: ReasonSessionNotFound, code: codes.NotFound}

func RaiseSessionNotFound(sessionID uuid.UUID, cause error) error {
	msg := fmt.Sprintf("session %s does not exist", sessionID)
	return &Error{
		code:    ErrSessionNotFound.code,
		reason:  ErrSessionNotFound.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "session",
				ResourceName: fmt.Sprintf("session:id:%s", sessionID),
				Description:  msg,
			},
		},
	}
}

var ErrSessionsForUserNotFound = &Error{reason: ReasonSessionsForUserNotFound, code: codes.NotFound}

func RaiseSessionsForUserNotFound(userID uuid.UUID, cause error) error {
	msg := fmt.Sprintf("no sessions found for user %s", userID)
	return &Error{
		code:    ErrSessionsForUserNotFound.code,
		reason:  ErrSessionsForUserNotFound.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "session",
				ResourceName: fmt.Sprintf("session:user_id:%s", userID),
				Description:  msg,
			},
		},
	}
}

var ErrSessionClientMismatch = &Error{reason: ReasonSessionClientMismatch, code: codes.PermissionDenied}

func RaiseSessionClientMismatch(sessionID uuid.UUID, cause error) error {
	msg := "session client mismatch"
	return &Error{
		code:    ErrSessionClientMismatch.code,
		reason:  ErrSessionClientMismatch.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "session",
				ResourceName: fmt.Sprintf("session:id:%s", sessionID),
				Description:  msg,
			},
			&errdetails.PreconditionFailure{
				Violations: []*errdetails.PreconditionFailure_Violation{
					{
						Type:        ErrSessionClientMismatch.reason,
						Subject:     sessionID.String(),
						Description: msg,
					},
				},
			},
		},
	}
}

var ErrSessionTokenMismatch = &Error{reason: ReasonSessionTokenMismatch, code: codes.Unauthenticated}

func RaiseSessionTokenMismatch(sessionID uuid.UUID, cause error) error {
	msg := "session token mismatch"
	return &Error{
		code:    ErrSessionTokenMismatch.code,
		reason:  ErrSessionTokenMismatch.reason,
		message: msg,
		cause:   cause,
		details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "session",
				ResourceName: fmt.Sprintf("session:id:%s", sessionID),
				Description:  msg,
			},
			&errdetails.PreconditionFailure{
				Violations: []*errdetails.PreconditionFailure_Violation{
					{
						Type:        ErrSessionTokenMismatch.reason,
						Subject:     fmt.Sprintf("session:id:%s", sessionID),
						Description: msg,
					},
				},
			},
		},
	}
}

var ErrSessionDoesNotBelongToUser = &Error{reason: ReasonSessionDoesNotBelongToUser, code: codes.PermissionDenied}

func RaiseSessionDoesNotBelongToUser(sessionID, userID uuid.UUID) error {
	msg := fmt.Sprintf("session %s does not belong to user %s", sessionID, userID)
	return &Error{
		code:    ErrSessionDoesNotBelongToUser.code,
		reason:  ErrSessionDoesNotBelongToUser.reason,
		message: msg,
		cause:   fmt.Errorf("session %s does not belong to user %s", sessionID, userID),
		details: []protoadapt.MessageV1{
			&errdetails.ResourceInfo{
				ResourceType: "session",
				ResourceName: fmt.Sprintf("session:id:%s", sessionID),
				Description:  msg,
			},
			&errdetails.PreconditionFailure{
				Violations: []*errdetails.PreconditionFailure_Violation{
					{
						Type:        ErrSessionDoesNotBelongToUser.reason,
						Subject:     fmt.Sprintf("session:id:%s", sessionID),
						Description: msg,
					},
				},
			},
		},
	}
}

var ErrNoPermission = &Error{reason: ReasonNoPermissions, code: codes.PermissionDenied}

func RaiseNoPermissions(cause error) error {
	msg := "no permission to perform this action"
	return &Error{
		code:    ErrNoPermission.code,
		reason:  ErrNoPermission.reason,
		message: msg,
		cause:   cause,
	}
}

var ErrUserSuspended = &Error{reason: ReasonUserSuspended, code: codes.FailedPrecondition}

func RaiseUserSuspended(userID uuid.UUID) error {
	msg := fmt.Sprintf("user %s is suspended", userID)
	return &Error{
		reason:  ErrUserSuspended.reason,
		message: msg,
		cause:   fmt.Errorf("user %s is suspended", userID),
		details: []protoadapt.MessageV1{
			&errdetails.PreconditionFailure{
				Violations: []*errdetails.PreconditionFailure_Violation{
					{
						Type:        ErrUserSuspended.reason,
						Subject:     fmt.Sprintf("user:id:%s/suspend", userID),
						Description: msg,
					},
				},
			},
		},
	}
}
