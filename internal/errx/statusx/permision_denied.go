package statusx

import (
	"time"

	"github.com/chains-lab/sso-svc/internal/constant"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	PermissionDeniedCode   = codes.PermissionDenied
	PermissionDeniedReason = "PERMISSION_DENIED"
)

func PermissionDenied(message string) *status.Status {
	response, _ := status.New(
		PermissionDeniedCode,
		message,
	).WithDetails(
		&errdetails.ErrorInfo{
			Reason: PermissionDeniedReason,
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	)

	return response
}
