package presenter

import (
	"net/http"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/google/uuid"
)

func (p Presenter) AppError(w http.ResponseWriter, requestID uuid.UUID, appErr *ape.Error) {
	errorID := uuid.New()
	p.log.WithField("request_id", requestID).
		WithField("error_id", errorID).
		WithField("code", appErr.Code).
		WithError(appErr.Unwrap()).
		Error("error from app")

	base := httpkit.ResponseErrorInput{
		Code:      appErr.Code,
		Title:     appErr.Title,
		Detail:    appErr.Details,
		RequestID: requestID.String(),
		ErrorID:   errorID.String(),
	}

	switch appErr.Code {
	// 404 Not Found
	case ape.CodeUserDoesNotExist,
		ape.CodeSessionDoesNotExist,
		ape.CodeSessionsForUserNotExist:
		base.Status = http.StatusNotFound

	// 409 Conflict
	case ape.CodeUserAlreadyExists,
		ape.CodeSessionAlreadyExists,
		ape.CodeSessionClientMismatch,
		ape.CodeSessionTokenMismatch:
		base.Status = http.StatusConflict

	// 400 Bad Request
	case ape.CodeUserInvalidRole,
		ape.CodeSessionCannotBeCurrent,
		ape.CodeSessionCannotBeCurrentUser,
		ape.CodeInvalidRequestBody,
		ape.CodeInvalidRequestQuery,
		ape.CodeInvalidRequestHeader,
		ape.CodeInvalidRequestPath:
		base.Status = http.StatusBadRequest

	// 403 Forbidden
	case ape.CodeUserNoPermissionToUpdateRole,
		ape.CodeSessionCannotDeleteSuperUserByOther:
		base.Status = http.StatusForbidden

	// 401 Unauthorized
	case ape.UnauthorizedError:
		base.Status = http.StatusUnauthorized

	// 500 Internal Server Error (fallback for any other code)
	default:
		base.Status = http.StatusInternalServerError
	}

	httpkit.RenderErr(w, httpkit.ResponseError(base)...)
}
