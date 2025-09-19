package user

import (
	"context"
	"fmt"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (u User) CompareRightsForAdmins(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	dif int,
) (initiator, user models.User, err error) {
	initiator, err = u.GetInitiator(ctx, initiatorID)
	if err != nil {
		return initiator, user, err
	}

	user, err = u.GetByID(ctx, userID)
	if err != nil {
		return initiator, user, err
	}

	if initiatorID == userID {
		return initiator, user, nil
	}

	if user.Role != roles.SuperUser {
		allowed, err := roles.CompareRolesUser(initiator.Role, user.Role)
		if err != nil {
			return initiator, user, errx.ErrorRoleNotSupported.Raise(
				fmt.Errorf("comparing roles between initiator %s and user %s, cause: %w", initiator.Role, user.Role, err),
			)
		}

		if allowed < dif {
			return initiator, user, errx.ErrorNoPermissions.Raise(
				fmt.Errorf("initiator Role %s is not allowed to interact with this user", initiator.Role),
			)
		}
	}

	return initiator, user, nil
}
