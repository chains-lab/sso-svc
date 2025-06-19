package handlers

import (
	"context"

	"github.com/chains-lab/proto-storage/gen/go/sso"
	"golang.org/x/oauth2"
)

func (h Handlers) GoogleLogin(ctx context.Context, request *sso.Empty) (*sso.GoogleLoginResponse, error) {
	url := h.google.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	// Вместо http.Redirect — возвращаем его в теле ответа
	return &sso.GoogleLoginResponse{Url: url}, nil
}
