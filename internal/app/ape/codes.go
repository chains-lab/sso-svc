package ape

const (
	//General error codes

	CodeInvalidRequestBody   = "INVALID_REQUEST_BODY"
	CodeInvalidRequestQuery  = "INVALID_REQUEST_QUERY"
	CodeInvalidRequestHeader = "INVALID_REQUEST_HEADER"
	CodeInvalidRequestPath   = "INVALID_REQUEST_PATH"
	CodeInvalidRequestMethod = "INVALID_REQUEST_METHOD"
	UnauthorizedError        = "UNAUTHORIZED"

	//Specific error codes

	CodeAccountDoesNotExist                 = "ACCOUNT_DOES_NOT_EXIST"
	CodeSessionDoesNotExist                 = "SESSION_DOES_NOT_EXIST"
	CodeAccountAlreadyExists                = "ACCOUNT_ALREADY_EXISTS"
	CodeAccountInvalidRole                  = "ACCOUNT_INVALID_ROLE"
	CodeUserNoPermissionToUpdateRole        = "USER_NO_PERMISSION_UPDATE_ROLE"
	CodeSessionsForAccountNotExist          = "SESSIONS_FOR_ACCOUNT_NOT_EXIST"
	CodeSessionClientMismatch               = "SESSIONS_CLIENT_MISMATCH"
	CodeSessionTokenMismatch                = "SESSIONS_TOKEN_MISMATCH"
	CodeSessionAlreadyExists                = "SESSION_ALREADY_EXISTS"
	CodeSessionCannotBeCurrent              = "SESSION_CANNOT_BE_CURRENT"
	CodeSessionCannotBeCurrentAccount       = "SESSION_CANNOT_BE_CURRENT_ACCOUNT"
	CodeSessionCannotDeleteSuperUserByOther = "SESSION_CANNOT_DELETE_SUPERUSER_BY_OTHER"
	CodeInternal                            = "INTERNAL_SERVER_ERROR"
)
