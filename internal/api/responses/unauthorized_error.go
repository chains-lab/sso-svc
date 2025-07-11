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

func UnauthorizedError(
	ctx context.Context,
	reason string,
	requestID *uuid.UUID,
) error {
	st := status.New(codes.Unauthenticated, "bad credentials")

	info := &errdetails.ErrorInfo{
		Reason: ape.ReasonUnauthorized,
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
