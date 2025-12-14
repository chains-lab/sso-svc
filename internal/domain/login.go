package domain

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) LoginByEmail(ctx context.Context, email, password string) (entity.TokensPair, error) {
	account, err := s.GetAccountByEmail(ctx, email)
	if err != nil {
		return entity.TokensPair{}, err
	}

	if err = account.CanInteract(); err != nil {
		return entity.TokensPair{}, err
	}

	err = s.CheckAccountPassword(ctx, account.ID, password)
	if err != nil {
		return entity.TokensPair{}, err
	}

	return s.CreateSession(ctx, account)
}

func (s Service) LoginByUsername(ctx context.Context, username, password string) (entity.TokensPair, error) {
	account, err := s.GetAccountByUsername(ctx, username)
	if err != nil {
		return entity.TokensPair{}, err
	}

	if err = account.CanInteract(); err != nil {
		return entity.TokensPair{}, err
	}

	err = s.CheckAccountPassword(ctx, account.ID, password)
	if err != nil {
		return entity.TokensPair{}, err
	}

	return s.CreateSession(ctx, account)
}

func (s Service) LoginByGoogle(ctx context.Context, email string) (entity.TokensPair, error) {
	account, err := s.GetAccountByEmail(ctx, email)
	if err != nil {
		return entity.TokensPair{}, err
	}

	if err = account.CanInteract(); err != nil {
		return entity.TokensPair{}, err
	}

	return s.CreateSession(ctx, account)
}

func (s Service) CheckAccountPassword(
	ctx context.Context,
	accountID uuid.UUID,
	password string,
) error {
	passData, err := s.db.GetAccountPassword(ctx, accountID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get account password, cause: %w", err),
		)
	}

	if err = passData.CheckPasswordMatch(password); err != nil {
		return err
	}

	return nil
}
