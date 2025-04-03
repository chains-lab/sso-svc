package data

import (
	"database/sql"
	"time"

	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/hs-zavet/sso-oauth/internal/service/infra/data/cache"
	"github.com/hs-zavet/sso-oauth/internal/service/infra/data/sqldb"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Data struct {
	SQL   SQLStorage
	Cache CacheStorage
}

type SQLStorage struct {
	Accounts sqldb.Accounts
	Sessions sqldb.Sessions
}

type CacheStorage struct {
	Accounts cache.Accounts
	Sessions cache.Sessions
}

func NewData(cfg *config.Config, log *logrus.Logger) (*Data, error) {
	db, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Database.Redis.Addr,
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.DB,
	})

	sqlAccounts := sqldb.NewAccounts(db)
	sqlSessions := sqldb.NewSessions(db)

	redisAccounts := cache.NewAccounts(redisClient, time.Duration(cfg.Database.Redis.Lifetime)*time.Minute)
	redisSessions := cache.NewSessions(redisClient, time.Duration(cfg.Database.Redis.Lifetime)*time.Minute)

	return &Data{
		SQL: SQLStorage{
			Accounts: sqlAccounts,
			Sessions: sqlSessions,
		},
		Cache: CacheStorage{
			Accounts: redisAccounts,
			Sessions: redisSessions,
		},
	}, nil
}
