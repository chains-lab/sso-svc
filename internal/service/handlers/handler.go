package handlers

import "github.com/recovery-flow/sso-oauth/internal/config"

type Handlers struct {
	svc *config.Service
}

func NewHandlers(svc *config.Service) *Handlers {
	return &Handlers{svc: svc}
}
