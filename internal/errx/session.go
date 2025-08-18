package errx

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/constant"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorSessionNotFound = ape.Declare("SESSION_NOT_FOUND")

func RaiseSessionNotFound(ctx context.Context, cause error, sessionID, userID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("session not found for user %s with session ID: %s", userID, sessionID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorSessionNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)

	return ErrorSessionNotFound.Raise(cause, st)
}

var ErrorSessionsForUserNotFound = ape.Declare("SESSIONS_FOR_USER_NOT_FOUND")

func RaiseSessionsForUserNotFound(ctx context.Context, cause error) error {
	st := status.New(codes.NotFound, fmt.Sprintf("sessions not found for user"))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorSessionsForUserNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)

	return ErrorSessionsForUserNotFound.Raise(cause, st)
}

var ErrorSessionTokenMismatch = ape.Declare("SESSION_TOKEN_MISMATCH")

func RaiseSessionTokenMismatch(ctx context.Context, cause error) error {
	st := status.New(codes.PermissionDenied, "session token mismatch")
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorSessionTokenMismatch.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)

	return ErrorSessionTokenMismatch.Raise(cause, st)
}

var ErrorSessionClientMismatch = ape.Declare("SESSION_CLIENT_MISMATCH")

func RaiseSessionClientMismatch(ctx context.Context, cause error) error {
	st := status.New(codes.FailedPrecondition, "session client mismatch")
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorSessionClientMismatch.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)

	return ErrorSessionClientMismatch.Raise(cause, st)
}
