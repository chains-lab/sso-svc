package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorUserNotFound = ape.DeclareError("USER_NOT_FOUND")

var ErrorUserAlreadyExists = ape.DeclareError("USER_ALREADY_EXISTS")

var ErrorLoginIsIncorrect = ape.DeclareError("LOGIN_IS_INCORRECT")
