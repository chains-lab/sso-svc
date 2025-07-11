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

type Violation struct {
	Field       string
	Description string
}

func BadRequestError(
	ctx context.Context,
	requestID uuid.UUID,
	violations ...Violation,
) error {
	st := status.New(codes.InvalidArgument, "bad request")

	info := &errdetails.ErrorInfo{
		Reason: ape.ReasonBadRequest,
		Domain: ape.ServiceName,
		Metadata: map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	var fb []*errdetails.BadRequest_FieldViolation
	for _, v := range violations {
		fb = append(fb, &errdetails.BadRequest_FieldViolation{
			Field:       v.Field,
			Description: v.Description,
		})
	}
	br := &errdetails.BadRequest{FieldViolations: fb}

	ri := &errdetails.RequestInfo{
		RequestId: requestID.String(),
	}

	st, err := st.WithDetails(info, br, ri)
	if err != nil {
		// если не удалось упаковать — возвращаем без деталей
		return st.Err()
	}

	return st.Err()
}
