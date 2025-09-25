package user

import (
	"github.com/chains-lab/sso-svc/internal/data"
)

type Service struct {
	db data.Database
}

func CreateUser(db data.Database) Service {
	return Service{
		db: db,
	}
}

func NewService(db data.Database) Service {
	return Service{
		db: db,
	}
}
