package service

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/app/ape"
	"github.com/google/uuid"
)

func (s Service) GetUserSessionsByAdmin(ctx context.Context, req *svc.GetUserSessionsByAdminRequest) (*svc.SessionsList, error) {
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &svc.SessionsList{}, responses.BadRequestError(ctx, meta.RequestID, ape.Violation{
			Field:       "user_id",
			Description: "invalid format user id",
		})
	}

	sessions, err := s.app.GetUserSessions(ctx, userID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Errorf("failed to retrieve sessions for user %s", req.UserId)

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Infof("retrieved sessions for user %s by admin %s", req.UserId, meta.InitiatorID)
	return responses.SessionList(sessions), nil
}
