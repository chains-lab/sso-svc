package statusx

import (
	"time"

	"github.com/chains-lab/sso-svc/internal/constant"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	InvalidArgumentCode   = codes.InvalidArgument
	InvalidArgumentReason = "INVALID_ARGUMENT"
)

func InvalidArgument(message, field, description string) *status.Status {
	response, _ := status.New(
		InvalidArgumentCode,
		message,
	).WithDetails(
		&errdetails.ErrorInfo{
			Reason: InvalidArgumentReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
		&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{{
				Field:       field,
				Description: description,
			}},
		},
	)

	return response
}
