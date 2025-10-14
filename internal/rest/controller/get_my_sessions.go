package controller

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/restkit/pagi"
	"github.com/chains-lab/sso-svc/internal/rest/meta"
	"github.com/chains-lab/sso-svc/internal/rest/responses"
)

func (s *Service) GetMySessions(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	page, size := pagi.GetPagination(r)
	sessions, err := s.domain.Session.ListForUser(r.Context(), initiator.ID, page, size)
	if err != nil {
		s.log.WithError(err).Errorf("failed to select My sessions")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.UserSessionsCollection(sessions))
}
