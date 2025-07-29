package ape

const (
	ReasonInternal                   = "INTERNAL_ERROR"
	ReasonBadRequest                 = "BAD_REQUEST"
	ReasonUnauthorized               = "UNAUTHORIZED"
	ReasonUserNotFound               = "USER_DOES_NOT_EXIST"
	ReasonSessionNotFound            = "SESSION_DOES_NOT_EXIST"
	ReasonUserAlreadyExists          = "USER_ALREADY_EXISTS"
	ReasonSessionsForUserNotFound    = "SESSIONS_FOR_USER_NOT_EXIST"
	ReasonSessionClientMismatch      = "SESSION_CLIENT_MISMATCH"
	ReasonSessionTokenMismatch       = "SESSION_TOKEN_MISMATCH"
	ReasonSessionDoesNotBelongToUser = "SESSION_DOES_NOT_BELONG_TO_USER"
	ReasonNoPermissions              = "NO_PERMISSIONS"
	ReasonUserSuspended              = "USER_SUSPENDED"
)
