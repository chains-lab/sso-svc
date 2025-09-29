package user

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (s Service) Delete(ctx context.Context, userID uuid.UUID) error {
	err := s.db.Users().FilterID(userID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(err)
	}

	return nil
}
