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
		switch appErr.Reason {
		case ape.ReasonUserDoesNotExist,
			ape.ReasonSessionDoesNotExist:

			code = codes.NotFound

		case ape.ReasonUserAlreadyExists,
			ape.ReasonSessionsForUserNotExist,
			ape.ReasonSessionClientMismatch,
			ape.ReasonSessionTokenMismatch:

			code = codes.AlreadyExists

		case ape.ReasonInternal:
			code = codes.Internal

		default:
			code = codes.Unknown
		}

		st := status.New(code, appErr.Reason)
		st, errWithDetails := st.WithDetails(&errdetails.ErrorInfo{
			Reason: appErr.Reason,
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
