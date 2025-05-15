package controllers

import (
	"net/http"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/google/uuid"
)

func (h Controller) ResultFromApp(w http.ResponseWriter, requestID uuid.UUID, appErr *ape.Error) {
	errorID := uuid.New()
	h.log.WithField("request_id", requestID).
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
	// resource not found
	case ape.CodeAccountDoesNotExist,
		ape.CodeSessionDoesNotExist,
		ape.CodeSessionsForAccountNotExist:
		base.Status = http.StatusNotFound

	// conflict / already exists
	case ape.CodeAccountAlreadyExists,
		ape.CodeSessionAlreadyExists,
		ape.CodeSessionClientMismatch,
		ape.CodeSessionTokenMismatch:
		base.Status = http.StatusConflict

	// bad request
	case ape.CodeAccountInvalidRole,
		ape.CodeSessionCannotBeCurrent,
		ape.CodeSessionCannotBeCurrentAccount:
		base.Status = http.StatusBadRequest

	// forbidden
	case ape.CodeUserNoPermissionToUpdateRole,
		ape.CodeSessionCannotDeleteSuperUserByOther:
		base.Status = http.StatusForbidden

	// internal
	case ape.CodeInternal:
		base.Status = http.StatusInternalServerError

	// catch-all
	default:
		base.Status = http.StatusInternalServerError
		base.Code = ape.CodeInternal
		base.Title = "Internal server error"
		base.Detail = "An unexpected error occurred"
	}

	httpkit.RenderErr(w, httpkit.ResponseError(base)...)
}
