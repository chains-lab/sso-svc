package data

import (
	"github.com/recovery-flow/sso-oauth/internal/data/dbx/redisdb"
	"github.com/recovery-flow/sso-oauth/internal/data/dbx/sql"
	"github.com/recovery-flow/sso-oauth/internal/data/repositories"
)

type Config struct {
	SqlUrl        string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

type Data struct {
	Accounts repositories.Accounts
	Sessions repositories.Sessions
}

func NewDataBase(cfg Config) (*Data, error) {
	queries, err := sql.NewRepoSQL(cfg.SqlUrl)
	if err != nil {
		return nil, err
	}
	redisDb, err := redisdb.NewRedisRepo(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		return nil, err
	}

	return &Data{
		Accounts: repositories.NewAccounts(redisDb.Accounts, queries.Accounts),
		Sessions: repositories.NewSessions(redisDb.Sessions, queries.Sessions),
	}, nil
}
