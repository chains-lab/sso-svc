package statusx

import (
	"time"

	"github.com/chains-lab/sso-svc/internal/constant"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AlreadyExistsCode   = codes.AlreadyExists
	AlreadyExistsReason = "ALREADY_EXISTS"
)

func AlreadyExists(message string) *status.Status {
	response, _ := status.New(
		AlreadyExistsCode,
		message,
	).WithDetails(
		&errdetails.ErrorInfo{
			Reason: AlreadyExistsReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}
