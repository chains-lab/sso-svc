package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/api/rest/requests"
	"github.com/chains-lab/sso-svc/internal/api/rest/responses"
)

func (s Service) Login(w http.ResponseWriter, r *http.Request) {
	req, err := requests.Login(r)
	if err != nil {
		s.Log(r).WithError(err).Error("failed to decode login request")

		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	token, err := s.app.Login(r.Context(), req.Data.Attributes.Email, req.Data.Attributes.Password, "TODO", "TODO")
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to login user")

		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	s.Log(r).Infof("user %s logged in successfully", req.Data.Attributes.Email)

	ape.Render(w, http.StatusOK, responses.TokensPair(token))
}
