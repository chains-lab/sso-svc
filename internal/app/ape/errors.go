package ape

import (
	"fmt"
)

var (
	ErrAccountNotFound                            = fmt.Errorf("account not found")
	ErrAccountAlreadyExists                       = fmt.Errorf("account already exists")
	ErrAccountInvalidRole                         = fmt.Errorf("invalid role")
	ErrSessionNotFound                            = fmt.Errorf("session not found")
	ErrSessionsNotFound                           = fmt.Errorf("sessions not found")
	ErrSessionsClientMismatch                     = fmt.Errorf("client does not match")
	ErrSessionsTokenMismatch                      = fmt.Errorf("token does not match")
	ErrSessionAlreadyExists                       = fmt.Errorf("session already exists")
	ErrSessionCannotBeCurrent                     = fmt.Errorf("session cannot be current")
	ErrSessionCannotBeCurrentAccount              = fmt.Errorf("session cannot be current account")
	ErrSessionCannotDeleteForSuperUserByOtherUser = fmt.Errorf("cannot delete superuser session by other user")
)
