package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/app/jwtmanager"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
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
	FilterRole(role string) dbx.UserQ
	FilterVerified(verified bool) dbx.UserQ

	Update(ctx context.Context, input dbx.UserUpdateInput) error

	Page(limit, offset uint64) dbx.UserQ
	Transaction(fn func(ctx context.Context) error) error
}

//type Broker interface {
//	AdminCreateUser(ctx context.Context, created bodies.UserCreated) error
//}

type Users struct {
	repo usersQ
	jwt  JWTManager
	//kafka Broker
}

func NewUser(cfg config.Config, log logger.Logger) (Users, error) {
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

type UserCreateInput struct {
	Email    string `json:"email"`
	Role     string `json:"role"`
	Verified bool   `json:"verified"`
}

func (a Users) Create(ctx context.Context, input UserCreateInput) error {
	ID := uuid.New()
	CreatedAt := time.Now().UTC()

	txErr := a.repo.New().Transaction(func(ctx context.Context) error {
		if err := a.repo.New().Insert(ctx, dbx.UserModel{
			ID:        ID,
			Email:     input.Email,
			Role:      input.Role,
			Verified:  input.Verified,
			UpdatedAt: CreatedAt,
			CreatedAt: CreatedAt,
		}); err != nil {
			return err
		}

		//TODO: in future we can use kafka to notify other services about user creation
		//if err := a.kafka.AdminCreateUser(ctx, bodies.UserCreated{
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
			return errx.RaiseUserAlreadyExists(ctx, txErr, input.Email)
		default:
			return errx.RaiseInternal(ctx, txErr)
		}
	}

	return nil
}

func (a Users) GetByID(ctx context.Context, ID uuid.UUID) (models.User, error) {
	user, err := a.repo.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.RaiseUserNotFound(
				ctx,
				fmt.Errorf("user with ID '%s' not found cause: %s", ID, err),
				ID.String(),
			)
		default:
			return models.User{}, errx.RaiseInternal(ctx, err)
		}
	}

	return models.User{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		Verified:  user.Verified,
		Suspended: user.Suspended,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (a Users) GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.User, error) {
	initiator, err := a.GetByID(ctx, initiatorID)
	if err != nil {
		return models.User{}, errx.RaiseInitiatorNotFound(
			ctx,
			fmt.Errorf("initiator with ID '%s' not found cause: %s", initiatorID, err),
			initiatorID,
		)
	}

	return initiator, nil
}

func (a Users) GetByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := a.repo.New().FilterEmail(email).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, errx.RaiseUserNotFound(
				ctx,
				fmt.Errorf("user with email '%s' not found cause: %s", email, err),
				email,
			)
		default:
			return models.User{}, errx.RaiseInternal(ctx, err)
		}
	}

	return models.User{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		Verified:  user.Verified,
		Suspended: user.Suspended,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (a Users) UpdateRole(ctx context.Context, ID uuid.UUID, role string) error {
	UpdatedAt := time.Now().UTC()

	err := a.repo.New().FilterID(ID).Update(ctx, dbx.UserUpdateInput{
		Role:      &role,
		UpdatedAt: UpdatedAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.RaiseUserNotFound(
				ctx,
				fmt.Errorf("user with ID '%s' not found cause: %s", ID, err),
				ID.String(),
			)
		default:
			return errx.RaiseInternal(ctx, err)
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
			return errx.RaiseUserNotFound(
				ctx,
				fmt.Errorf("user with ID '%s' not found cause: %s", ID, err),
				ID.String(),
			)
		default:
			return errx.RaiseInternal(ctx, err)
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
			return errx.RaiseUserNotFound(
				ctx,
				fmt.Errorf("user with ID '%s' not found cause: %s", ID, err),
				ID.String(),
			)
		default:
			return errx.RaiseInternal(ctx, err)
		}
	}

	return nil
}

func (a Users) ComparisonRightsForAdmins(ctx context.Context, initiatorID, userID uuid.UUID) (initiator, user models.User, err error) {
	initiator, err = a.GetByID(ctx, initiatorID)
	if err != nil {
		return initiator, user, err
	}

	user, err = a.GetByID(ctx, userID)
	if err != nil {
		return initiator, user, err
	}

	if initiator.Suspended {
		return initiator, user, errx.RaiseInitiatorUserSuspended(
			ctx,
			fmt.Errorf("initiator %s is suspended", initiatorID),
			initiatorID.String(),
		)
	}

	if user.Role == roles.User {
		return initiator, user, errx.RaiseUserRoleIsNotAllowed(
			ctx,
			fmt.Errorf("initiator Role %s is not allowed to interact wit with this user", initiator.Role),
		)
	}

	if initiatorID == userID {
		return initiator, user, nil
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, user.Role) < 1 {
			return initiator, user, errx.RaiseInitiatorRoleIsLowThanTarget(
				ctx,
				fmt.Errorf("initiator Role %s is not allowed to interact with this user", initiator.Role),
			)
		}
	}

	return initiator, user, nil
}
