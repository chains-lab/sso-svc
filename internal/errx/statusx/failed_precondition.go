package statusx

import (
	"time"

	"github.com/chains-lab/sso-svc/internal/constant"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	FailedPreconditionCode   = codes.FailedPrecondition
	FailedPreconditionReason = "FAILED_PRECONDITION"
)

func FailedPrecondition(message string) *status.Status {
	response, _ := status.New(
		FailedPreconditionCode,
		message,
	).WithDetails(
		&errdetails.ErrorInfo{
			Reason: FailedPreconditionReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}
