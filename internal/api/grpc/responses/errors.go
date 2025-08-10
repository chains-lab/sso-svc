package responses

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/constant"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AppError(
	ctx context.Context,
	err error,
) error {
	requestID := meta.RequestID(ctx).String()
	if requestID == uuid.Nil.String() {
		requestID = "unknown"
	}

	st, ok := status.FromError(err)
	if !ok {
		return InternalError(ctx)
	}

	withReq, derr := st.WithDetails(
		&errdetails.RequestInfo{RequestId: requestID},
	)
	if derr != nil {
		return status.Errorf(
			codes.Internal,
			"failed to attach request info: %v",
			derr,
		)
	}

	return withReq.Err()
}

func InternalError(
	ctx context.Context,
) error {
	requestID := meta.RequestID(ctx).String()
	if requestID == uuid.Nil.String() {
		requestID = "unknown"
	}

	st := status.New(codes.Internal, "internal server error")

	info := &errdetails.ErrorInfo{
		Reason: "INTERNAL",
		Domain: constant.ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	ri := &errdetails.RequestInfo{
		RequestId: requestID,
	}

	st, err := st.WithDetails(info, ri)
	if err != nil {
		return st.Err()
	}

	return st.Err()
}

func InvalidArgumentError(
	ctx context.Context,
	message string,
	violations ...*errdetails.BadRequest_FieldViolation,
) error {
	requestID := meta.RequestID(ctx).String()
	if requestID == uuid.Nil.String() {
		requestID = "unknown"
	}

	st := status.New(codes.InvalidArgument, message)

	info := &errdetails.ErrorInfo{
		Reason: "INVALID_ARGUMENT",
		Domain: constant.ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	br := &errdetails.BadRequest{FieldViolations: violations}

	ri := &errdetails.RequestInfo{
		RequestId: requestID,
	}

	st, err := st.WithDetails(info, br, ri)
	if err != nil {
		return st.Err()
	}

	return st.Err()
}

func UnauthenticatedError(
	ctx context.Context,
	message string,
) error {
	requestID := meta.RequestID(ctx).String()
	if requestID == uuid.Nil.String() {
		requestID = "unknown"
	}

	st := status.New(codes.Unauthenticated, message)

	info := &errdetails.ErrorInfo{
		Reason: "UNAUTHENTICATED",
		Domain: constant.ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	ri := &errdetails.RequestInfo{
		RequestId: requestID,
	}

	st, err := st.WithDetails(info, ri)
	if err != nil {
		return st.Err()
	}

	return st.Err()
}

func PermissionDeniedError(
	ctx context.Context,
	message string,
) error {
	requestID := meta.RequestID(ctx).String()
	if requestID == uuid.Nil.String() {
		requestID = "unknown"
	}

	st := status.New(codes.PermissionDenied, message)

	info := &errdetails.ErrorInfo{
		Reason: "PERMISSION_DENIED",
		Domain: constant.ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	ri := &errdetails.RequestInfo{
		RequestId: requestID,
	}

	st, err := st.WithDetails(info, ri)
	if err != nil {
		return st.Err()
	}

	return st.Err()
}
