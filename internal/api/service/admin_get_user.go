package service

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/app/ape"
	"github.com/google/uuid"
)

func (s Service) GetUserByAdmin(ctx context.Context, req *svc.GetUserByAdminRequest) (*svc.User, error) {
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to parse user ID")

		return &svc.User{}, responses.BadRequestError(ctx, meta.RequestID, ape.Violation{
			Field:       "user_id",
			Description: "invalid format user id",
		})
	}

	user, err := s.app.GetUserByID(ctx, userID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to get user")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return &svc.User{
		Id:           user.ID.String(),
		Email:        user.Email,
		Role:         string(user.Role),
		Subscription: user.Subscription.String(),
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}
