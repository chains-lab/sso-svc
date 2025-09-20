package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/app/domain/session"
	"github.com/chains-lab/sso-svc/internal/app/domain/user"
	"github.com/chains-lab/sso-svc/internal/app/jwtmanager"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/google/uuid"
)

type App struct {
	session session.Session
	users   user.User

	db *sql.DB
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	manager := jwtmanager.NewManager(cfg)

	return App{
		session: session.CreateSession(pg, manager),
		users:   user.CreateUser(pg),

		db: pg,
	}, nil
}

func (a App) transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, dbx.TxKey, tx)

	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (a App) getInitiator(ctx context.Context, userID, sessionID uuid.UUID) (models.User, error) {
	_, err := a.session.GetSessionForInitiator(ctx, userID, sessionID)
	if err != nil {
		return models.User{}, err
	}

	return a.users.GetByID(ctx, userID)
}
