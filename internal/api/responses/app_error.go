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
	var appErr *ape.Error
	if errors.As(err, &appErr) {

		st := status.New(appErr.Code(), appErr.Error())

		details := appErr.Details()
		details = append(details, &errdetails.RequestInfo{
			RequestId: requestID.String(),
		})
		if err != nil {
			return st.Err()
		}
	}

	return status.Errorf(codes.Internal, "unexpected error")
}
