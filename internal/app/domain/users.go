package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/chains-lab/chains-auth/internal/events/bodies"
	"github.com/chains-lab/chains-auth/internal/events/writer"
	"github.com/chains-lab/chains-auth/internal/jwtkit"
	"github.com/chains-lab/chains-auth/internal/repo"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type usersRepo interface {
	Create(ctx context.Context, input repo.UserModel) error
	Update(ctx context.Context, ID uuid.UUID, input repo.UserUpdateRequest) error
	Delete(ctx context.Context, ID uuid.UUID) error
	GetByID(ctx context.Context, ID uuid.UUID) (repo.UserModel, error)
	GetByEmail(ctx context.Context, email string) (repo.UserModel, error)
	Transaction(fn func(ctx context.Context) error) error
	Drop(ctx context.Context) error
}

type Broker interface {
	CreateUser(ctx context.Context, created bodies.UserCreated) error
}

type Users struct {
	repo  usersRepo
	jwt   JWTManager
	kafka Broker
}

func NewUser(cfg config.Config, log *logrus.Logger) (Users, error) {
	data, err := repo.NewUsers(cfg, log)
	if err != nil {
		return Users{}, nil
	}

	jwt := jwtkit.NewManager(cfg)

	kafka := writer.NewUserCreateWriters(cfg, log.WithFields(logrus.Fields{
		"component": "kafka",
		"topic":     bodies.UserCreateTopic,
	}))

	return Users{
		repo:  data,
		jwt:   jwt,
		kafka: kafka,
	}, nil
}

func (a Users) Create(ctx context.Context, email string, role roles.Role) error {
	ID := uuid.New()
	CreatedAt := time.Now().UTC()

	txErr := a.repo.Transaction(func(ctx context.Context) error {
		if err := a.repo.Create(ctx, repo.UserModel{
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

// TODO maybe is doesn't need
func (a Users) UpdateRole(ctx context.Context, ID uuid.UUID, role roles.Role) error {
	UpdatedAt := time.Now().UTC()

	err := a.repo.Update(ctx, ID, repo.UserUpdateRequest{
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

func (a Users) GetByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	user, err := a.repo.GetByID(ctx, ID)
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
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (a Users) GetByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := a.repo.GetByEmail(ctx, email)
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
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}
