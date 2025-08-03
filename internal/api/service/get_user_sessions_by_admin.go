package service

import (
	"context"

	"github.com/chains-lab/sso-proto/gen/go/svc"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) GetUserSessionsByAdmin(ctx context.Context, req *svc.GetUserSessionsByAdminRequest) (*svc.SessionsList, error) {
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &svc.SessionsList{}, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid format user id",
		})
	}

	sessions, err := s.app.GetUserSessions(ctx, userID, req.Pagination.Page, req.Pagination.Limit)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Errorf("failed to retrieve sessions for user %s", req.UserId)

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("retrieved sessions for user %s by admin %s", req.UserId, meta.InitiatorID)
	return responses.SessionList(sessions, req.Pagination.Page), nil
}
