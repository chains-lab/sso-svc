package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) UpdateCompany(
	ctx context.Context,
	userID uuid.UUID,
	companyID *uuid.UUID,
	role *string,
) error {
	err := s.db.Transaction(ctx, func(txCtx context.Context) error {
		err := s.db.UpdateUserCompany(ctx, userID, companyID, role, time.Now().UTC())
		if err != nil {
			return errx.ErrorInternal.Raise(fmt.Errorf("failed to update user comapny: %w", err))
		}

		err = s.db.DeleteAllSessionsForUser(txCtx, userID)
		if err != nil {
			return errx.ErrorInternal.Raise(fmt.Errorf("failed to delete user sessions after city update: %w", err))
		}

		return nil
	})

	return err
}
