package presenter

import (
	"net/http"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/google/uuid"
)

func (p Presenter) InvalidToken(w http.ResponseWriter, requestID uuid.UUID, err error) {
	errorID := uuid.New()

	p.log.WithField("request_id", requestID).
		WithField("error_id", errorID).
		WithError(err).
		Error("error getting user data from token")

	httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
		Status:    http.StatusUnauthorized,
		Code:      ape.UnauthorizedError,
		Detail:    err.Error(),
		RequestID: requestID.String(),
		ErrorID:   errorID.String(),
	})...)
}
