package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/chains-lab/chains-auth/internal/repo/redisdb"
	"github.com/chains-lab/chains-auth/internal/repo/sqldb"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Role         roles.Role `json:"role"`
	Subscription uuid.UUID  `json:"subscription,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type usersSQL interface {
	New() sqldb.UserQ
	Insert(ctx context.Context, input sqldb.UserInsertInput) error
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int, error)
	Select(ctx context.Context) ([]sqldb.UserModel, error)
	Get(ctx context.Context) (sqldb.UserModel, error)

	FilterID(id uuid.UUID) sqldb.UserQ
	FilterEmail(email string) sqldb.UserQ
	FilterRole(role roles.Role) sqldb.UserQ
	FilterSubscription(subscription uuid.UUID) sqldb.UserQ

	Update(ctx context.Context, input sqldb.UserUpdateInput) error

	Page(limit, offset uint64) sqldb.UserQ
	Transaction(fn func(ctx context.Context) error) error

	Drop(ctx context.Context) error
}

type usersRedis interface {
	Create(ctx context.Context, input redisdb.CreateUserInput) error
	Set(ctx context.Context, input redisdb.UserSetInput) error
	Update(ctx context.Context, userID uuid.UUID, input redisdb.UserUpdateRequest) error
	GetByID(ctx context.Context, userID string) (redisdb.UserModel, error)
	GetByEmail(ctx context.Context, email string) (redisdb.UserModel, error)
	Delete(ctx context.Context, userID, email string) error

	Drop(ctx context.Context) error
}

type UsersRepo struct {
	sql   usersSQL
	redis usersRedis
	log   *logrus.Entry
}

func NewUsers(cfg config.Config, log *logrus.Logger) (UsersRepo, error) {
	db, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return UsersRepo{}, err
	}

	redisImpl := redisdb.NewUsers(redis.NewClient(&redis.Options{
		Addr:     cfg.Database.Redis.Addr,
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.DB,
	}), cfg.Database.Redis.Lifetime)
	sqlImpl := sqldb.NewUsers(db)

	return UsersRepo{
		sql:   sqlImpl,
		redis: redisImpl,
		log:   log.WithField("component", "users"),
	}, nil
}

type UserCreateRequest struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Role         roles.Role `json:"role"`
	Subscription uuid.UUID  `json:"subscription,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (a UsersRepo) Create(ctx context.Context, input UserCreateRequest) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	if err := a.sql.Insert(ctxSync, sqldb.UserInsertInput{
		ID:           input.ID,
		Email:        input.Email,
		Role:         input.Role,
		Subscription: input.Subscription,
		CreatedAt:    input.CreatedAt,
	}); err != nil {
		return err
	}

	_ = a.redis.Create(ctxSync, redisdb.CreateUserInput{
		ID:           input.ID,
		Email:        input.Email,
		Role:         input.Role,
		Subscription: input.Subscription,
	})

	return nil
}

type UserUpdateRequest struct {
	Role         *roles.Role `json:"role"`
	Subscription *uuid.UUID  `json:"subscription,omitempty"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

func (a UsersRepo) Update(ctx context.Context, ID uuid.UUID, input UserUpdateRequest) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	var sqlInput sqldb.UserUpdateInput
	if input.Role != nil {
		sqlInput.Role = input.Role
	}
	if input.Subscription != nil {
		sqlInput.Subscription = input.Subscription
	}
	sqlInput.UpdatedAt = input.UpdatedAt

	if err := a.sql.New().FilterID(ID).Update(ctxSync, sqldb.UserUpdateInput{
		Role:         input.Role,
		Subscription: input.Subscription,
		UpdatedAt:    input.UpdatedAt,
	}); err != nil {
		return err
	}

	user, err := a.sql.New().FilterID(ID).Get(ctxSync)
	if err != nil {
		return err
	}

	_ = a.redis.Set(ctxSync, redisdb.UserSetInput{
		ID:           user.ID,
		Email:        user.Email,
		Role:         user.Role,
		Subscription: user.Subscription,
		UpdatedAt:    user.UpdatedAt,
		CreatedAt:    user.CreatedAt,
	})
	return nil
}

func (a UsersRepo) Delete(ctx context.Context, ID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	user, err := a.sql.New().FilterID(ID).Get(ctxSync)
	if err != nil {
		return err
	}

	if err := a.redis.Delete(ctxSync, user.ID.String(), user.Email); err != nil {
		a.log.WithField("database", "redis").Errorf("error creating user in redis: %v", err)
	}

	if err := a.sql.New().FilterID(ID).Delete(ctxSync); err != nil {
		return err
	}

	return nil
}

func (a UsersRepo) GetByID(ctx context.Context, ID uuid.UUID) (User, error) {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	redisRes, err := a.redis.GetByID(ctxSync, ID.String())
	if err != nil {
		a.log.WithField("database", "redis").Errorf("error creating user in redis: %v", err)
	} else {
		res := User{
			ID:           redisRes.ID,
			Email:        redisRes.Email,
			Subscription: redisRes.Subscription,
			CreatedAt:    redisRes.CreatedAt,
			Role:         redisRes.Role,
		}
		if redisRes.UpdatedAt != nil {
			res.UpdatedAt = *redisRes.UpdatedAt
		}
		return res, nil
	}

	user, err := a.sql.New().FilterID(ID).Get(ctxSync)
	if err != nil {
		return User{}, err
	}

	res := User{
		ID:           user.ID,
		Email:        user.Email,
		Role:         user.Role,
		Subscription: user.Subscription,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    *user.UpdatedAt,
	}
	if user.UpdatedAt != nil {
		res.UpdatedAt = *user.UpdatedAt
	}

	if err := a.redis.Set(ctxSync, redisdb.UserSetInput{
		ID:           user.ID,
		Email:        user.Email,
		Role:         user.Role,
		Subscription: user.Subscription,
		UpdatedAt:    user.UpdatedAt,
		CreatedAt:    user.CreatedAt,
	}); err != nil {
		a.log.WithField("database", "redis").Errorf("error creating user in redis: %v", err)
	}

	return res, nil
}

func (a UsersRepo) GetByEmail(ctx context.Context, email string) (User, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	redisRes, err := a.redis.GetByEmail(ctxSync, email)
	if err != nil {
		a.log.WithField("database", "redis").Errorf("error creating user in redis: %v", err)
	} else {
		res := User{
			ID:           redisRes.ID,
			Email:        redisRes.Email,
			Subscription: redisRes.Subscription,
			CreatedAt:    redisRes.CreatedAt,
			Role:         redisRes.Role,
		}
		if redisRes.UpdatedAt != nil {
			res.UpdatedAt = *redisRes.UpdatedAt
		}
		return res, nil
	}

	user, err := a.sql.New().FilterEmail(email).Get(ctxSync)
	if err != nil {
		return User{}, err
	}

	res := User{
		ID:           user.ID,
		Email:        user.Email,
		Role:         user.Role,
		Subscription: user.Subscription,
		CreatedAt:    user.CreatedAt,
	}

	if user.UpdatedAt != nil {
		res.UpdatedAt = *user.UpdatedAt
	}

	if err := a.redis.Set(ctxSync, redisdb.UserSetInput{
		ID:           user.ID,
		Email:        user.Email,
		Role:         user.Role,
		Subscription: user.Subscription,
		UpdatedAt:    user.UpdatedAt,
		CreatedAt:    user.CreatedAt,
	}); err != nil {
		a.log.WithField("database", "redis").Errorf("error creating user in redis: %v", err)
	}

	return res, nil
}

func (a UsersRepo) Transaction(fn func(ctx context.Context) error) error {
	return a.sql.Transaction(fn)
}

func (a UsersRepo) Drop(ctx context.Context) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	if err := a.sql.Drop(ctxSync); err != nil {
		return err
	}

	if err := a.redis.Drop(ctxSync); err != nil {
		return err
	}

	return nil
}
