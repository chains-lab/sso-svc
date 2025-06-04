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

	CodeUserDoesNotExist                    = "ACCOUNT_DOES_NOT_EXIST"
	CodeSessionDoesNotExist                 = "SESSION_DOES_NOT_EXIST"
	CodeUserAlreadyExists                   = "ACCOUNT_ALREADY_EXISTS"
	CodeUserInvalidRole                     = "ACCOUNT_INVALID_ROLE"
	CodeUserNoPermissionToUpdateRole        = "USER_NO_PERMISSION_UPDATE_ROLE"
	CodeSessionsForUserNotExist             = "SESSIONS_FOR_ACCOUNT_NOT_EXIST"
	CodeSessionClientMismatch               = "SESSIONS_CLIENT_MISMATCH"
	CodeSessionTokenMismatch                = "SESSIONS_TOKEN_MISMATCH"
	CodeSessionAlreadyExists                = "SESSION_ALREADY_EXISTS"
	CodeSessionCannotBeCurrent              = "SESSION_CANNOT_BE_CURRENT"
	CodeSessionCannotBeCurrentUser          = "SESSION_CANNOT_BE_CURRENT_ACCOUNT"
	CodeSessionCannotDeleteSuperUserByOther = "SESSION_CANNOT_DELETE_SUPERUSER_BY_OTHER"
)
