package user

import (
	"github.com/chains-lab/sso-svc/internal/data"
)

type Service struct {
	db data.Database
}

func New(db data.Database) Service {
	return Service{
		db: db,
	}
}
