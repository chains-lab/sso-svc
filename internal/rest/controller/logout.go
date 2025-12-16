package controller

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/domain/modules/auth"
	"github.com/chains-lab/sso-svc/internal/rest/meta"
)

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.AccountData(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	err = s.domain.Logout(r.Context(), auth.InitiatorData{
		AccountID: initiator.ID,
		SessionID: initiator.SessionID,
	})
	if err != nil {
		s.log.WithError(err).Errorf("failed to logout user")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
