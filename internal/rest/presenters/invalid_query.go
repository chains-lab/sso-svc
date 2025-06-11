package presenters

import (
	"fmt"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/google/uuid"
)

func (p Presenters) InvalidQuery(w http.ResponseWriter, requestID uuid.UUID, query string, err error) {
	errorID := uuid.New()

	p.log.WithField("request_id", requestID).
		WithField("error_id", errorID).
		WithError(err).
		Errorf("invalid query %s parameter", query)

	httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
		Status: http.StatusBadRequest,
		Code:   ape.CodeInvalidRequestQuery,
		Error:  err,
		Title:  "Invalid query parameter",
		Detail: fmt.Sprintf("The query parameter '%s' is invalid: %s", query, err.Error()),
	})...)
}
