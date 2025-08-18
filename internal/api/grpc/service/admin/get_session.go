package admin

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/svc/admin"
	sesionProto "github.com/chains-lab/sso-proto/gen/go/svc/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/logger"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetSession(ctx context.Context, req *svc.GetSessionRequest) (*sesionProto.Session, error) {
	//TODO implement authorization check for admin role
	initiator, err := meta.User(ctx)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	sessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, problems.InvalidArgumentError(ctx, "invalid session ID format", &errdetails.BadRequest_FieldViolation{
			Field:       "session_id",
			Description: "invalid UUID format for session ID",
		})
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, problems.InvalidArgumentError(ctx, "invalid user ID format", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	session, err := s.app.AdminGetUserSession(ctx, initiator.ID, userID, sessionID)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("failed to get session for user %s with session ID %s", req.UserId, req.SessionId)

		return nil, err
	}

	return response.Session(session), nil
}
