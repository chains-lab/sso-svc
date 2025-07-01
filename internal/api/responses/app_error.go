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
	var appErr *ape.BusinessError
	if errors.As(err, &appErr) {
		var code codes.Code
		switch appErr.Reason() {
		case ape.ReasonUserDoesNotExist,
			ape.ReasonSessionDoesNotExist,
			ape.ReasonSessionsForUserNotExist:

			code = codes.NotFound

		case ape.ReasonUserAlreadyExists:

			code = codes.AlreadyExists

		case ape.ReasonSessionClientMismatch,
			ape.ReasonSessionTokenMismatch:

			code = codes.Unauthenticated

		case ape.ReasonSessionDoesNotBelongToUser,
			ape.ReasonNoPermission,
			ape.ReasonUserSuspended:

			code = codes.PermissionDenied

		case ape.ReasonInternal:
			code = codes.Internal

		default:
			code = codes.Unknown
		}

		st := status.New(code, appErr.Error())

		info := &errdetails.ErrorInfo{
			Reason: appErr.Reason(),
			Domain: "sso-svc",
			Metadata: map[string]string{
				"request_id": requestID.String(),
			},
		}

		if code == codes.InvalidArgument {
			var fb []*errdetails.BadRequest_FieldViolation

			for _, v := range appErr.Violations() {
				fb = append(fb, &errdetails.BadRequest_FieldViolation{
					Field:       v.Field,
					Description: v.Description,
				})
			}
			br := &errdetails.BadRequest{FieldViolations: fb}

			st, err := st.WithDetails(info, br)
			if err != nil {
				return st.Err()
			}
		}

		return st.Err()
	}

	return status.Errorf(codes.Internal, "Unexcpected error")
}
