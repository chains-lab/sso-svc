package guard

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
)

func AllowedRoles(ctx context.Context, req *userdata.UserData, action string, allowed ...string) (uuid.UUID, error) {
	initiatorID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse initiator ID")

		return uuid.Nil, problems.UnauthenticatedError(ctx, "invalid initiator ID format")
	}

	allow := false
	for _, role := range allowed {
		if req.Role == role {
			allow = true
			break
		}
	}

	if !allow {
		logger.Log(ctx).Warnf(
			"user %s with role %s tried to perform this action: '%s', that requires one of the allowed roles: %v",
			req.UserId, req.Role, action, allowed,
		)

		return uuid.Nil, problems.PermissionDeniedError(ctx,
			fmt.Sprintf("initiator role can perform this '%s' action", action))
	}
	return initiatorID, nil
}
