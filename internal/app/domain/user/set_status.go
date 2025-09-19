package user

import (
	"context"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (u User) SetStatus(ctx context.Context, userID uuid.UUID, status string) error {
	err := enum.ParseUserStatus(status)
	if err != nil {
		return errx.ErrorUserStatusNotSupported.Raise(
			fmt.Errorf("parsing status for user %s, cause: %w", userID, err),
		)
	}

	err = u.query.New().FilterID(userID).Update(ctx,
		map[string]interface{}{"status": status},
	)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("updating status for user %s, cause: %w", userID, err),
		)
	}

	return nil
}
