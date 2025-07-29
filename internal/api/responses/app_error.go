package responses

import (
	"context"
	"errors"
	"time"

	"github.com/chains-lab/sso-svc/internal/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

func AppError(ctx context.Context, requestID uuid.UUID, err error) error {
	var appErr *ape.Error
	if errors.As(err, &appErr) {

		st := status.New(appErr.Code(), appErr.Error())

		details := appErr.Details()
		details = append(details, &errdetails.ErrorInfo{
			Reason: appErr.Reason(),
			Domain: ape.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		})
		details = append(details, &errdetails.RequestInfo{
			RequestId: requestID.String(),
		})
		if err != nil {
			return st.Err()
		}
	}

	return InternalError(ctx, &requestID)
}
