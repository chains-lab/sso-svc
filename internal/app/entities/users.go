package entities

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/app/ape"
	"github.com/chains-lab/sso-svc/internal/app/jwtmanager"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type usersQ interface {
	New() dbx.UserQ
	Insert(ctx context.Context, input dbx.UserModel) error
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int, error)
	Select(ctx context.Context) ([]dbx.UserModel, error)
	Get(ctx context.Context) (dbx.UserModel, error)

	FilterID(id uuid.UUID) dbx.UserQ
	FilterEmail(email string) dbx.UserQ
	FilterRole(role roles.Role) dbx.UserQ
	FilterSubscription(subscription uuid.UUID) dbx.UserQ
	FilterVerified(verified bool) dbx.UserQ

	Update(ctx context.Context, input dbx.UserUpdateInput) error

	Page(limit, offset uint64) dbx.UserQ
	Transaction(fn func(ctx context.Context) error) error

	//Drop(ctx context.Context) error
}

//type Broker interface {
//	CreateUser(ctx context.Context, created bodies.UserCreated) error
//}

type Users struct {
	repo usersQ
	jwt  JWTManager
	//kafka Broker
}

func NewUser(cfg config.Config, log *logrus.Logger) (Users, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return Users{}, err
	}

	//kafka := writer.NewUserCreateWriters(cfg, log.WithFields(logrus.Fields{
	//	"component": "kafka",
	//	"topic":     bodies.UserCreateTopic,
	//}))

	return Users{
		repo: dbx.NewUsers(pg),
		jwt:  jwtmanager.NewManager(cfg),
		//kafka: kafka,
	}, nil
}

func (a Users) Create(ctx context.Context, email string, role roles.Role) error {
	ID := uuid.New()
	CreatedAt := time.Now().UTC()

	txErr := a.repo.New().Transaction(func(ctx context.Context) error {
		if err := a.repo.New().Insert(ctx, dbx.UserModel{
			ID:           ID,
			Email:        email,
			Role:         role,
			Subscription: uuid.Nil,
			Verified:     false,
			UpdatedAt:    CreatedAt,
			CreatedAt:    CreatedAt,
		}); err != nil {
			return err
		}

		//TODO: in future we can use kafka to notify other services about user creation
		//if err := a.kafka.CreateUser(ctx, bodies.UserCreated{
		//	UserID:    ID.String(),
		//	Role:      role,
		//	Timestamp: CreatedAt,
		//}); err != nil {
		//	return err
		//}

		return nil
	})
	if txErr != nil {
		switch {
		case errors.Is(txErr, sql.ErrNoRows):
			return ape.ErrorUserAlreadyExists(txErr)
		default:
			return ape.ErrorInternal(txErr)
		}
	}

	return nil
}

func (a Users) GetByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	user, err := a.repo.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, ape.ErrorUserDoesNotExist(ID, err)
		default:
			return models.User{}, ape.ErrorInternal(err)
		}
	}

	return models.User{
		ID:           user.ID,
		Email:        user.Email,
		Role:         user.Role,
		Subscription: user.Subscription,
		Verified:     user.Verified,
		Suspended:    user.Suspended,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (a Users) GetByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := a.repo.New().FilterEmail(email).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, ape.ErrorUserDoesNotExistByEmail(email, err)
		default:
			return models.User{}, ape.ErrorInternal(err)
		}
	}

	return models.User{
		ID:           user.ID,
		Email:        user.Email,
		Role:         user.Role,
		Subscription: user.Subscription,
		Verified:     user.Verified,
		Suspended:    user.Suspended,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (a Users) UpdateRole(ctx context.Context, ID uuid.UUID, role roles.Role) error {
	UpdatedAt := time.Now().UTC()

	err := a.repo.New().FilterID(ID).Update(ctx, dbx.UserUpdateInput{
		Role:      &role,
		UpdatedAt: UpdatedAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorUserDoesNotExist(ID, err)
		default:
			return ape.ErrorInternal(err)
		}
	}

	return nil
}

func (a Users) UpdateSubscription(ctx context.Context, ID, subscriptionID uuid.UUID) error {
	UpdatedAt := time.Now().UTC()

	err := a.repo.New().FilterID(ID).Update(ctx, dbx.UserUpdateInput{
		Subscription: &subscriptionID,
		UpdatedAt:    UpdatedAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorUserDoesNotExist(ID, err)
		default:
			return ape.ErrorInternal(err)
		}
	}

	return nil
}

func (a Users) UpdateVerified(ctx context.Context, ID uuid.UUID, verified bool) error {
	UpdatedAt := time.Now().UTC()

	err := a.repo.New().FilterID(ID).Update(ctx, dbx.UserUpdateInput{
		Verified:  &verified,
		UpdatedAt: UpdatedAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorUserDoesNotExist(ID, err)
		default:
			return ape.ErrorInternal(err)
		}
	}

	return nil
}

func (a Users) UpdateSuspended(ctx context.Context, ID uuid.UUID, suspended bool) error {
	UpdatedAt := time.Now().UTC()

	err := a.repo.New().FilterID(ID).Update(ctx, dbx.UserUpdateInput{
		Suspended: &suspended,
		UpdatedAt: UpdatedAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorUserDoesNotExist(ID, err)
		default:
			return ape.ErrorInternal(err)
		}
	}

	return nil
}
