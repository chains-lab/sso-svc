package presenters

import (
	"fmt"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/google/uuid"
)

func (p Presenters) MismatchIdentification(w http.ResponseWriter, requestID uuid.UUID, parameter, pointer string) {
	ErrorID := uuid.New()

	p.log.WithField("request_id", requestID).
		WithField("error_id", ErrorID).
		WithError(fmt.Errorf("invalid URL or JSON resource ID"))

	httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
		Status:    http.StatusBadRequest,
		Code:      ape.CodeInvalidRequestPath,
		Title:     "Invalid URL or JSON resource ID",
		Detail:    "Invalid URL or JSON resource ID",
		ErrorID:   ErrorID.String(),
		RequestID: requestID.String(),
		Parameter: parameter,
		Pointer:   pointer,
	})...)

	return
}
