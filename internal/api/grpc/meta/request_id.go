package meta

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/api/grpc/interceptor"
	"github.com/google/uuid"
)

func RequestID(ctx context.Context) string {
	if ctx == nil {
		return "unknow"
	}

	requestID, ok := ctx.Value(interceptor.RequestIDCtxKey).(uuid.UUID)
	if !ok {
		return "unknow"
	}

	return requestID.String()
}
