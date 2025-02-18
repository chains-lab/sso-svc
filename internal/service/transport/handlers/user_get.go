package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/responses"
	"github.com/recovery-flow/tokens"
)

func (h *Handler) AccountGet(w http.ResponseWriter, r *http.Request) {
	svc := h.svc
	log := svc.Logger

	userID, ok := r.Context().Value(tokens.UserIDKey).(uuid.UUID)
	if !ok {
		log.Warn("UserID not found in context")
		httpkit.RenderErr(w, problems.Unauthorized("Accounts not authenticated"))
		return
	}

	user, err := svc.DB.Accounts.GetByID(r, userID)
	if err != nil {
		log.Errorf("Failed to retrieve user: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
	httpkit.Render(w, responses.Account(*user))
}
