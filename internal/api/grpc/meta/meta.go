package meta

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/api/grpc/interceptor"
	"github.com/google/uuid"
)

func RequestID(ctx context.Context) uuid.UUID {
	if ctx == nil {
		return uuid.Nil
	}

	requestID, ok := ctx.Value(interceptor.RequestIDCtxKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return requestID
}
