package ape

import (
	"fmt"
)

var (
	ErrAccountDoseNotExits  = fmt.Errorf("account does not exist")
	ErrAccountAlreadyExists = fmt.Errorf("account already exists")

	ErrAccountInvalidRole = fmt.Errorf("invalid role")

	ErrUserHasNoPermissionToUpdateRole = fmt.Errorf("user has no permission to update role")

	ErrSessionDoseNotExits                        = fmt.Errorf("session doses not exist")
	ErrSessionsForAccountDoseNotExits             = fmt.Errorf("sessions for account doses not exist")
	ErrSessionsClientMismatch                     = fmt.Errorf("client does not match")
	ErrSessionsTokenMismatch                      = fmt.Errorf("token does not match")
	ErrSessionAlreadyExists                       = fmt.Errorf("session already exists")
	ErrSessionCannotBeCurrent                     = fmt.Errorf("session cannot be current")
	ErrSessionCannotBeCurrentAccount              = fmt.Errorf("session cannot be current account")
	ErrSessionCannotDeleteForSuperUserByOtherUser = fmt.Errorf("cannot delete superuser session by other user")
)
