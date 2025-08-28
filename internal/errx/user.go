package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorUserNotFound = ape.DeclareError("USER_NOT_FOUND")

var ErrorUserAlreadyExists = ape.DeclareError("USER_ALREADY_EXISTS")

var ErrorLoginIsIncorrect = ape.DeclareError("LOGIN_IS_INCORRECT")

var ErrorInitiatorIsBlocked = ape.DeclareError("USER_IS_BLOCKED")

var ErrorRoleNotSupported = ape.DeclareError("USER_ROLE_NOT_SUPPORTED")

var ErrorUserStatusNotSupported = ape.DeclareError("USER_STATUS_NOT_SUPPORTED")
