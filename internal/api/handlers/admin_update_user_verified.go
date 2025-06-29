package handlers

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Service) AdminUpdateUserVerified(ctx context.Context, req *svc.AdminUpdateUserVerifiedRequest) (*svc.UserResponse, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	user, err := s.app.AdminUpdateUserVerified(ctx, meta.InitiatorID, userID, req.Verified)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Warnf("user %s verified status updated to %t by %s", req.UserId, req.Verified, meta.InitiatorID)
	return responses.User(user), nil
}
