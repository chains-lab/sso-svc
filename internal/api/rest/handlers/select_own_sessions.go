package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/api/rest/meta"
	"github.com/chains-lab/sso-svc/internal/api/rest/responses"
)

func (s Service) SelectOwnSessions(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	pagReq, sort := pagi.GetPagination(r)

	sessions, pag, err := s.app.ListOwnSessions(r.Context(), initiator.UserID, pagReq, sort)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to select own sessions")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.UserSessionsCollection(sessions, pag))
}
