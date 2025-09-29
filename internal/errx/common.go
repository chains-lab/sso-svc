package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorInternal = ape.DeclareError("INTERNAL_ERROR")

var ErrorNoPermissions = ape.DeclareError("NO_PERMISSIONS")

var ErrorUnauthenticated = ape.DeclareError("UNAUTHENTICATED")
