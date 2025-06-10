package presenter

import (
	"net/http"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/google/uuid"
)

func (p Presenter) InvalidPointer(w http.ResponseWriter, requestID uuid.UUID, err error) {
	errorID := uuid.New()

	p.log.WithField("request_id", requestID).
		WithField("error_id", errorID).
		WithError(err).
		Errorf("error processing request with invalid pointer")

	httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
		Status: http.StatusBadRequest,
		Code:   ape.CodeInvalidRequestBody,
		Error:  err,
	})...)
	return
}
