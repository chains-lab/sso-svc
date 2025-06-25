package handlers

import (
	"context"

	"github.com/chains-lab/proto-storage/gen/go/auth"
	"golang.org/x/oauth2"
)

func (a Service) GoogleLogin(ctx context.Context, request *auth.Empty) (*auth.GoogleLoginResponse, error) {
	url := a.cfg.GoogleOAuth().AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	// Вместо http.Redirect — возвращаем его в теле ответа
	return &auth.GoogleLoginResponse{Url: url}, nil
}
