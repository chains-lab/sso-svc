package db

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/cifra-city/comtools/httpkit"
	"github.com/cifra-city/sso-oauth/internal/data/db/sqlcore"
	"github.com/google/uuid"
)

type Transaction interface {
	LoginTxn(
		r *http.Request,
		userID uuid.UUID,
		deviceName string,
		Token string,
		deviceID uuid.UUID,
	) (*sqlcore.Session, error)

	TerminateSessionsTxn(
		r *http.Request,
		userId uuid.UUID,
		curDevId uuid.UUID,
	) error

	UpdateRefreshTokenTrx( //TODO for future use right no sense
		r *http.Request,
		userID uuid.UUID,
		sessionID uuid.UUID,
		newToken string,
	) error
}

type transaction struct {
	queries *sqlcore.Queries
}

func NewTransaction(queries *sqlcore.Queries) Transaction {
	return &transaction{queries: queries}
}

func HandleTransactionRollback(tx *sql.Tx, originalErr error) error {
	if originalErr != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", originalErr, rbErr)
		}
	}
	return originalErr
}

func (t *transaction) LoginTxn(
	r *http.Request,
	userID uuid.UUID,
	deviceName string,
	Token string,
	deviceID uuid.UUID,
) (*sqlcore.Session, error) {
	ctx := r.Context()
	queries, tx, err := t.queries.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = HandleTransactionRollback(tx, err)
	}()

	client := httpkit.GetUserAgent(r)

	if err != nil {
		return nil, err
	}

	session, err := queries.CreateSession(ctx, sqlcore.CreateSessionParams{
		ID:     deviceID,
		UserID: userID,
		Token:  Token,
		Client: client,
	})
	if err != nil {
		return nil, err
	}

	return &session, tx.Commit()
}

func (t *transaction) TerminateSessionsTxn(
	r *http.Request,
	userId uuid.UUID,
	curDevId uuid.UUID,
) error {
	ctx := r.Context()
	queries, tx, err := t.queries.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = HandleTransactionRollback(tx, err)
	}()

	if err != nil {
		return err
	}

	userSessions, err := queries.GetSessionsByUserID(ctx, userId)
	if err != nil {
		return err
	}

	for _, dev := range userSessions {
		if dev.ID == curDevId {
			continue
		}
		err = queries.DeleteUserSession(ctx, sqlcore.DeleteUserSessionParams{
			ID:     dev.ID,
			UserID: userId,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (t *transaction) UpdateRefreshTokenTrx( //TODO for future use
	r *http.Request,
	userID uuid.UUID,
	sessionID uuid.UUID,
	newToken string) error {

	ctx := r.Context()
	queries, tx, err := t.queries.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = HandleTransactionRollback(tx, err)
	}()

	err = queries.UpdateSessionToken(ctx, sqlcore.UpdateSessionTokenParams{
		ID:     sessionID,
		UserID: userID,
		Token:  newToken,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}
