package responses

import (
	"context"
	"errors"

	"github.com/chains-lab/sso-svc/internal/app/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AppError(ctx context.Context, requestID uuid.UUID, err error) error {
	errorID := uuid.New()
	var appErr *ape.Error
	if errors.As(err, &appErr) {
		var code codes.Code
		switch {
		case errors.Is(appErr.Err, ape.ErrUserDoesNotExist),
			errors.Is(appErr.Err, ape.ErrSessionDoesNotExist):

			code = codes.NotFound

		case errors.Is(appErr.Err, ape.ErrUserAlreadyExists),
			errors.Is(appErr.Err, ape.ErrSessionsForUserNotExist),
			errors.Is(appErr.Err, ape.ErrSessionClientMismatch),
			errors.Is(appErr.Err, ape.ErrSessionTokenMismatch):

			code = codes.AlreadyExists

		case errors.Is(appErr.Err, ape.ErrInternal):
			code = codes.Internal

		default:
			code = codes.Unknown
		}

		st := status.New(code, appErr.Err.Error())
		st, errWithDetails := st.WithDetails(&errdetails.ErrorInfo{
			Reason: appErr.Err.Error(),
			Metadata: map[string]string{
				"error_id":   errorID.String(),
				"request_id": requestID.String(),
			},
		})
		if errWithDetails != nil {
			return st.Err()
		}

		return st.Err()
	}

	return status.Errorf(codes.Internal, "Unexcpected error")
}
