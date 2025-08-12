package sessionadmin

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	sesionProto "github.com/chains-lab/sso-proto/gen/go/svc/session"
	svc "github.com/chains-lab/sso-proto/gen/go/svc/sessionadmin"
	"github.com/chains-lab/sso-svc/internal/api/grpc/guard"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problem"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetSession(ctx context.Context, req *svc.GetSessionRequest) (*sesionProto.Session, error) {
	_, err := guard.AllowedRoles(ctx, req.Initiator, "get user session by admin",
		roles.Admin, roles.SuperUser)
	if err != nil {
		return nil, err
	}

	_, user, err := s.ComparisonRightsForAdmins(ctx, req.Initiator.UserId, req.UserId)
	if err != nil {
		return nil, err
	}

	sessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, problem.InvalidArgumentError(ctx, "invalid session ID format", &errdetails.BadRequest_FieldViolation{
			Field:       "session_id",
			Description: "invalid UUID format for session ID",
		})
	}

	session, err := s.app.GetUserSession(ctx, user.ID, sessionID)
	if err != nil {
		return nil, err
	}

	return response.Session(session), nil
}
