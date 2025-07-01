package service

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) GoogleLogin(ctx context.Context, _ *emptypb.Empty) (*svc.GoogleLoginResponse, error) {
	url := s.cfg.GoogleOAuth().AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	// Вместо http.Redirect — возвращаем его в теле ответа
	return &svc.GoogleLoginResponse{Url: url}, nil
}
