package errx

import (
	"github.com/chains-lab/sso-svc/internal/errx/statusx"
	"github.com/chains-lab/svc-errors/ape"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorInternal = ape.Declare("INTERNAL_ERROR")

func RaiseInternal(cause error) error {
	return ErrorInternal.Raise(
		cause,
		status.New(codes.Internal, "internal server error"),
	)
}

var ErrorNoPermissions = ape.Declare("NO_PERMISSIONS")

func RaiseNoPermissions(cause error) error {
	return ErrorNoPermissions.Raise(
		cause,
		statusx.PermissionDenied(cause.Error()),
	)
}
