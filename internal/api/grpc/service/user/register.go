package user

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/svc/user"
	"github.com/chains-lab/sso-svc/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) Register(ctx context.Context, req *svc.RegisterRequest) (*emptypb.Empty, error) {
	err := s.app.Register(ctx, req.Email, req.Password)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("error registering user with email %s", req.Email)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
