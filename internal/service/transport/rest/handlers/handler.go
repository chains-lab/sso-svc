package handlers

import (
	"github.com/recovery-flow/sso-oauth/internal/service/transport"
)

type Handlers struct {
	svc *transport.Transport
}

func NewHandlers(svc *transport.Transport) *Handlers {
	return &Handlers{svc: svc}
}
