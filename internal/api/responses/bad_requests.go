package responses

import (
	"context"

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
		Reason: "BAD_REQUEST",
		Domain: "sso-svc", // ваше API/сервис
		Metadata: map[string]string{
			"request_id": requestID.String(),
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

	st, err := st.WithDetails(info, br)
	if err != nil {
		// если не удалось упаковать — возвращаем без деталей
		return st.Err()
	}

	return st.Err()
}
