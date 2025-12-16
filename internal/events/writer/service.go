package writer

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/events/contracts"
	//"github.com/chains-lab/sso-svc/internal/events/repository"
)

type Service struct {
	repo repository
	addr string
}

type repository interface {
	CreateOutboxEvent(
		ctx context.Context,
		event contracts.Event,
	) error
}

func New(addr string, repo repository) *Service {
	return &Service{
		repo: repo,
		addr: addr,
	}
}
