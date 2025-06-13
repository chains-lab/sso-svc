package ape

const (
	//General error codes

	CodeInternal             = "INTERNAL_SERVER_ERROR"
	CodeInvalidRequestBody   = "INVALID_REQUEST_BODY"
	CodeInvalidRequestQuery  = "INVALID_REQUEST_QUERY"
	CodeInvalidRequestHeader = "INVALID_REQUEST_HEADER"
	CodeInvalidRequestPath   = "INVALID_REQUEST_PATH"
	UnauthorizedError        = "UNAUTHORIZED"

	//Specific error codes

	CodeUserDoesNotExist                    = "USER_DOES_NOT_EXIST"
	CodeSessionDoesNotExist                 = "SESSION_DOES_NOT_EXIST"
	CodeUserAlreadyExists                   = "USER_ALREADY_EXISTS"
	CodeUserInvalidRole                     = "USER_INVALID_ROLE"
	CodeUserNoPermissionToUpdateRole        = "USER_NO_PERMISSION_UPDATE_ROLE"
	CodeSessionsForUserNotExist             = "SESSIONS_FOR_USER_NOT_EXIST"
	CodeSessionClientMismatch               = "SESSIONS_CLIENT_MISMATCH"
	CodeSessionTokenMismatch                = "SESSIONS_TOKEN_MISMATCH"
	CodeSessionAlreadyExists                = "SESSION_ALREADY_EXISTS"
	CodeSessionCannotBeCurrent              = "SESSION_CANNOT_BE_CURRENT"
	CodeSessionCannotBeCurrentUser          = "SESSION_CANNOT_BE_CURRENT_USER"
	CodeSessionCannotDeleteSuperUserByOther = "SESSION_CANNOT_DELETE_SUPERUSER_BY_OTHER"
)
