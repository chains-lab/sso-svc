package service

import (
	"context"

	"github.com/chains-lab/sso-proto/gen/go/svc"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) GetUserByAdmin(ctx context.Context, req *svc.GetUserByAdminRequest) (*svc.User, error) {
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to parse user ID")

		return &svc.User{}, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid format user id",
		})
	}

	user, err := s.app.GetUserByID(ctx, userID)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to get user")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return &svc.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}
