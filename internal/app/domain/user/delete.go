package user

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (u User) Delete(ctx context.Context, userID uuid.UUID) error {
	err := u.passQ.New().FilterID(userID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(err)
	}

	err = u.query.New().FilterID(userID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(err)
	}

	return nil
}
