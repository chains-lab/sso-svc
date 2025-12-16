package domain

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
)

func (s Service) ValidateSession(
	ctx context.Context,
	initiator InitiatorData,
) (entity.Account, entity.Session, error) {
	account, err := s.db.GetAccountByID(ctx, initiator.AccountID)
	if err != nil {
		return entity.Account{}, entity.Session{}, errx.ErrorInitiatorNotFound.Raise(
			fmt.Errorf("failed to get account with id '%s', cause: %w", initiator.SessionID, err),
		)
	}
	if account.IsNil() {
		return entity.Account{}, entity.Session{}, errx.ErrorInitiatorNotFound.Raise(
			fmt.Errorf("account with id '%s' not found", initiator.SessionID),
		)
	}

	if err = account.CanInteract(); err != nil {
		return entity.Account{}, entity.Session{}, errx.ErrorInitiatorIsNotActive.Raise(
			fmt.Errorf("account with id '%s' cannot interact, cause: %w", initiator.AccountID, err),
		)
	}

	session, err := s.db.GetSession(ctx, initiator.SessionID)
	if err != nil {
		return entity.Account{}, entity.Session{}, errx.ErrorInitiatorInvalidSession.Raise(
			fmt.Errorf("failed to get session with id '%s', cause: %w", initiator.SessionID, err),
		)
	}
	if session.IsNil() || session.AccountID != initiator.AccountID {
		return entity.Account{}, entity.Session{}, errx.ErrorInitiatorInvalidSession.Raise(
			fmt.Errorf("session with id '%s' not found for account '%s'", initiator.SessionID, initiator.AccountID),
		)
	}

	return account, session, nil
}

func (s Service) GetInitiatorEmail(ctx context.Context, initiator InitiatorData) (entity.AccountEmail, error) {
	_, _, err := s.ValidateSession(ctx, initiator)
	if err != nil {
		return entity.AccountEmail{}, err
	}

	accountEmail, err := s.db.GetAccountEmail(ctx, initiator.AccountID)
	if err != nil {
		return entity.AccountEmail{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get account email repo with id '%s', cause: %w", initiator.AccountID, err),
		)
	}

	if accountEmail.IsNil() {
		return entity.AccountEmail{}, errx.ErrorInitiatorNotFound.Raise(
			fmt.Errorf("account email repo with id '%s' not found", initiator.AccountID),
		)
	}

	return accountEmail, nil
}
