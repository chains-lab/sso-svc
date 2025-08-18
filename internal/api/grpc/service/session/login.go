package session

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/svc/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/logger"
)

func (s Service) Login(
	ctx context.Context,
	req *svc.LoginRequest,
) (*svc.TokensPair, error) {
	_, tokensPair, err := s.app.Login(ctx, req.Email, req.Password, req.Client, req.Ip)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("error logging in user with email %s", req.Email)
		return nil, err
	}

	return response.TokensPair(tokensPair), nil
}
