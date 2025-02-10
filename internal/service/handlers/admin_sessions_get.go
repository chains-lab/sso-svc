package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/cifractx"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/responses"
	"github.com/sirupsen/logrus"
)

func AdminSessionsGet(w http.ResponseWriter, r *http.Request) {
	Server, err := cifractx.GetValue[*config.Server](r.Context(), config.SERVER)
	if err != nil {
		logrus.Errorf("Failed to retrieve service configuration %s", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
	log := Server.Logger

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		log.Errorf("Failed to parse user_id: %v", err)
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"user_id": validation.NewError("user_id", "Invalid user_id"),
		})...)
		return
	}

	sessions, err := Server.SqlDB.Sessions.GetSessions(r, userID)
	if err != nil {
		log.Errorf("Failed to retrieve user session: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
