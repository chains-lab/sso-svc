package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/recovery-flow/sso-oauth/internal/data/dbx/redisdb/repositories"
	"github.com/redis/go-redis/v9"
)

type Repo struct {
	Accounts repositories.Accounts
	Sessions repositories.Sessions
}

func NewRedisRepo(redisAddr, redisPassword string, redisDB int) (*Repo, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error redis: %v", err))
	}

	return &Repo{
		Accounts: repositories.NewRedisUserStore(client),
		Sessions: repositories.NewRedisSessionStore(client),
	}, nil
}
