package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorSessionNotFound = ape.DeclareError("SESSION_NOT_FOUND")

var ErrorSessionsForUserNotFound = ape.DeclareError("SESSIONS_FOR_USER_NOT_FOUND")

var ErrorSessionTokenMismatch = ape.DeclareError("SESSION_TOKEN_MISMATCH")

var ErrorSessionClientMismatch = ape.DeclareError("SESSION_CLIENT_MISMATCH")
