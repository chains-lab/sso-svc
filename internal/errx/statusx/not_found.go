package statusx

import (
	"time"

	"github.com/chains-lab/sso-svc/internal/constant"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	NotFoundCode   = codes.NotFound
	NotFoundReason = "NOT_FOUND"
)

func NotFound(message string) *status.Status {
	response, _ := status.New(
		NotFoundCode,
		message,
	).WithDetails(
		&errdetails.ErrorInfo{
			Reason: NotFoundReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}
