package responses

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InternalError(
	ctx context.Context,
	requestID *uuid.UUID,
) error {
	st := status.New(codes.Internal, "internal server error")

	info := &errdetails.ErrorInfo{
		Reason: ape.ReasonInternal,
		Domain: ape.ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	var ri *errdetails.RequestInfo
	if requestID != nil {
		ri = &errdetails.RequestInfo{
			RequestId: requestID.String(),
		}
	}

	st, err := st.WithDetails(info, ri)

	if err != nil {
		// если не удалось упаковать — возвращаем без деталей
		return st.Err()
	}

	return st.Err()
}
