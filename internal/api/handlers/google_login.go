package handlers

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"golang.org/x/oauth2"
)

func (a Service) GoogleLogin(ctx context.Context, request *svc.Empty) (*svc.GoogleLoginResponse, error) {
	url := a.cfg.GoogleOAuth().AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	// Вместо http.Redirect — возвращаем его в теле ответа
	return &svc.GoogleLoginResponse{Url: url}, nil
}
