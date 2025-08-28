package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/api/rest/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Service) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid user id: %s", chi.URLParam(r, "user_id"))

		ape.RenderErr(w, problems.InvalidParameter("user_id", err))
		return
	}

	user, err := s.app.GetUserByID(r.Context(), userID)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to get user by id: %s", userID)

		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.User(user))
}
