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

var ErrorUserNotFound = ape.Declare("USER_NOT_FOUND")

func RaiseUserNotFound(ctx context.Context, cause error, userID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("user with id: %s not found", userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorUserNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)
	return ErrorUserNotFound.Raise(cause, st)
}

func RaiseUserNotFoundByEmail(ctx context.Context, cause error, email string) error {
	st := status.New(codes.NotFound, fmt.Sprintf("user with email: %s not found", email))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorUserNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)
	return ErrorUserNotFound.Raise(cause, st)
}

var ErrorUserAlreadyExists = ape.Declare("USER_ALREADY_EXISTS")

func RaiseUserAlreadyExists(ctx context.Context, cause error, email string) error {
	st := status.New(codes.AlreadyExists, fmt.Sprintf("user with email: %s already exists", email))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorUserAlreadyExists.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)
	return ErrorUserAlreadyExists.Raise(cause, st)
}

//var ErrorUserRoleIsNotAllowed = ape.Declare("USER_ROLE_IS_NOT_ALLOWED")
//
//func RaiseUserRoleIsNotAllowed(ctx context.Context, cause error) error {
//	st := status.New(codes.PermissionDenied, cause.Error())
//	st, _ = st.WithDetails(
//		&errdetails.ErrorInfo{
//			Reason: ErrorUserRoleIsNotAllowed.Error(),
//			Domain: constant.ServiceName,
//			Metadata: map[string]string{
//				"timestamp": nowRFC3339Nano(),
//			},
//		},
//		&errdetails.RequestInfo{
//			RequestId: meta.RequestID(ctx),
//		},
//	)
//	return ErrorUserRoleIsNotAllowed.Raise(cause, st)
//}
//
//var ErrorInitiatorRoleIsLowThanTarget = ape.Declare("INITIATOR_ROLE_IS_LOWER_THAN_TARGET")
//
//func RaiseInitiatorRoleIsLowThanTarget(ctx context.Context, cause error) error {
//	st := status.New(codes.PermissionDenied, cause.Error())
//	st, _ = st.WithDetails(
//		&errdetails.ErrorInfo{
//			Reason: ErrorInitiatorRoleIsLowThanTarget.Error(),
//			Domain: constant.ServiceName,
//			Metadata: map[string]string{
//				"timestamp": nowRFC3339Nano(),
//			},
//		},
//		&errdetails.RequestInfo{
//			RequestId: meta.RequestID(ctx),
//		},
//	)
//	return ErrorInitiatorRoleIsLowThanTarget.Raise(cause, st)
//}

var ErrorLoginIsIncorrect = ape.Declare("LOGIN_IS_INCORRECT")

func RaiseLoginIsIncorrect(ctx context.Context, cause error) error {
	st := status.New(codes.PermissionDenied, "login is incorrect")
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorLoginIsIncorrect.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx),
		},
	)

	return ErrorLoginIsIncorrect.Raise(cause, st)
}
