package user

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (u Service) Delete(ctx context.Context, userID uuid.UUID) error {
	err := u.db.UsersPassword().FilterID(userID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(err)
	}

	err = u.db.Users().FilterID(userID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(err)
	}

	return nil
}
