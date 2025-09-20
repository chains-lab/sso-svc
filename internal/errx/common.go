package errx

import (
	"time"

	"github.com/chains-lab/ape"
)

func nowRFC3339Nano() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

var ErrorInternal = ape.DeclareError("INTERNAL_ERROR")

var ErrorNoPermissions = ape.DeclareError("NO_PERMISSIONS")

var ErrorUnauthenticated = ape.DeclareError("UNAUTHENTICATED")

var ErrorUserCannotBlockHimself = ape.DeclareError("USER_CANNOT_BLOCK_HIMSELF")
