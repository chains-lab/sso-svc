package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/responses"
)

func (h *Handlers) AdminSessionsGet(w http.ResponseWriter, r *http.Request) {
	svc := h.svc
	log := svc.Logger

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		log.Errorf("Failed to parse user_id: %v", err)
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"user_id": validation.NewError("user_id", "Invalid user_id"),
		})...)
		return
	}

	sessions, err := svc.DB.Sessions.SelectByUserID(r, userID)
	if err != nil {
		log.Errorf("Failed to retrieve user session: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
