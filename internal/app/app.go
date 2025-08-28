package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/app/entities"
	"github.com/chains-lab/sso-svc/internal/app/jwtmanager"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/google/uuid"
)

type sessionsQ interface {
	New() dbx.SessionsQ
	Insert(ctx context.Context, input dbx.Session) error
	Update(ctx context.Context, input map[string]any) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]dbx.Session, error)
	Get(ctx context.Context) (dbx.Session, error)

	FilterID(id uuid.UUID) dbx.SessionsQ
	FilterUserID(userID uuid.UUID) dbx.SessionsQ

	Page(limit, offset uint64) dbx.SessionsQ
	Count(ctx context.Context) (uint64, error)

	OrderCreatedAt(ascending bool) dbx.SessionsQ

	Transaction(fn func(ctx context.Context) error) error
}

type usersQ interface {
	New() dbx.UserQ
	Insert(ctx context.Context, input dbx.UserModel) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]dbx.UserModel, error)
	Get(ctx context.Context) (dbx.UserModel, error)

	FilterID(id uuid.UUID) dbx.UserQ
	FilterEmail(email string) dbx.UserQ
	FilterRole(role string) dbx.UserQ
	FilterEmailVer(verified bool) dbx.UserQ

	Update(ctx context.Context, input map[string]any) error

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.UserQ

	Transaction(fn func(ctx context.Context) error) error
}

type passQ interface {
	New() dbx.UserPassQ
	Insert(ctx context.Context, input dbx.UserPasswordModel) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]dbx.UserPasswordModel, error)
	Get(ctx context.Context) (dbx.UserPasswordModel, error)

	FilterID(id uuid.UUID) dbx.UserPassQ

	Update(ctx context.Context, input map[string]any) error

	Page(limit, offset uint64) dbx.UserPassQ
	Count(ctx context.Context) (uint64, error)
}

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	GenerateAccess(
		userID uuid.UUID,
		sessionID uuid.UUID,
		idn string,
	) (string, error)

	GenerateRefresh(
		userID uuid.UUID,
		sessionID uuid.UUID,
		idn string,
	) (string, error)
}

type App struct {
	session entities.Session
	users   entities.User

	jwt JWTManager

	db *sql.DB
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	manager := jwtmanager.NewManager(cfg)

	return App{
		session: entities.CreateSession(pg, manager),
		users:   entities.CreateUser(pg),
		jwt:     manager,
	}, nil
}

func (a App) Transaction(fn func(ctx context.Context) error) error {
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
