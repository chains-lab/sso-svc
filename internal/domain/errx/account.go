package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorInitiatorIsNotActive = ape.DeclareError("INITIATOR_IS_NOT_ACTIVE")
var ErrorInitiatorNotFound = ape.DeclareError("INITIATOR_NOT_FOUND")

var ErrorEmailAlreadyExist = ape.DeclareError("EMAIL_ALREADY_EXIST")
var ErrorEmailNotVerified = ape.DeclareError("EMAIL_NOT_VERIFIED")
var ErrorCannotChangeEmailYet = ape.DeclareError("CANNOT_CHANGE_EMAIL_YET")

var ErrorPasswordInvalid = ape.DeclareError("PASSWORD_INVALID")
var ErrorPasswordIsNotAllowed = ape.DeclareError("PASSWORD_IS_NOT_ALLOWED")
var ErrorCannotChangePasswordYet = ape.DeclareError("CANNOT_CHANGE_PASSWORD_YET")

var ErrorUsernameIsNotAllowed = ape.DeclareError("USERNAME_IS_NOT_ALLOWED")
var ErrorUsernameAlreadyTaken = ape.DeclareError("USERNAME_ALREADY_TAKEN")
var ErrorCannotChangeUsernameYet = ape.DeclareError("CANNOT_CHANGE_USERNAME_YET")

var ErrorRoleNotSupported = ape.DeclareError("ACCOUNT_ROLE_NOT_SUPPORTED")
