package handlers

import (
	"errors"
	"net/http"

	"github.com/recovery-flow/comtools/cifractx"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/requests"
	"github.com/sirupsen/logrus"
)

type OperationType string

const (
	RESET_PASSWORD  OperationType = "reset_password"
	CHANGE_PASSWIRD OperationType = "change_password"
	CHANGE_EMAIL    OperationType = "change_email"
	REGISTRATION    OperationType = "registration"
	LOGIN           OperationType = "login"
)

func (op OperationType) IsValid() bool {
	switch op {
	case RESET_PASSWORD, REGISTRATION, LOGIN, CHANGE_PASSWIRD, CHANGE_EMAIL:
		return true
	default:
		return false
	}
}
