package handlers

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handlers) GetUserSessions(ctx context.Context, req *sso.UserRequest) (*sso.SessionsListResponse, error) {
	requestID := uuid.New()

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id: %s", req.UserId)
	}

	sessions, err := h.app.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, h.presenter.AppError(requestID, err)
	}

	return responses.SessionList(sessions), nil
}
