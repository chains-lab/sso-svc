package ape

const (
	ReasonInternal                   = "INTERNAL_ERROR"
	ReasonUserDoesNotExist           = "USER_DOES_NOT_EXIST"
	ReasonSessionDoesNotExist        = "SESSION_DOES_NOT_EXIST"
	ReasonUserAlreadyExists          = "USER_ALREADY_EXISTS"
	ReasonSessionsForUserNotExist    = "SESSIONS_FOR_USER_NOT_EXIST"
	ReasonSessionClientMismatch      = "SESSION_CLIENT_MISMATCH"
	ReasonSessionTokenMismatch       = "SESSION_TOKEN_MISMATCH"
	ReasonSessionDoesNotBelongToUser = "SESSION_DOES_NOT_BELONG_TO_USER"
	ReasonNoPermission               = "NO_PERMISSION"
	ReasonUserSuspended              = "USER_SUSPENDED"
)
