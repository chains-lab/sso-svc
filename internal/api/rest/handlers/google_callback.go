package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/sso-svc/internal/api/rest/responses"
	"github.com/chains-lab/sso-svc/internal/errx"
)

func (s Service) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		ape.RenderErr(w, problems.InvalidParameter(
			"code", fmt.Errorf("code is required")),
		)

		return
	}

	token, err := s.google.Exchange(r.Context(), code)
	if err != nil {
		s.Log(r).WithError(err).Errorf("error exchanging code for user id: %s", code)
		ape.RenderErr(w, problems.InternalError())

		return
	}

	client := s.google.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		s.Log(r).WithError(err).Errorf("error getting user info from google")
		ape.RenderErr(w, problems.InternalError())

		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			s.Log(r).WithError(err).Errorf("error closing response body")
			ape.RenderErr(w, problems.InternalError())

			return
		}
	}(resp.Body)

	var userInfo struct {
		Email string `json:"email"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		s.Log(r).WithError(err).Errorf("error decoding user info from google")
		ape.RenderErr(w, problems.InternalError())

		return
	}

	tokensPair, err := s.app.LoginByGoogle(r.Context(), userInfo.Email)
	if err != nil {
		s.Log(r).WithError(err).Errorf("error logging in user: %s", userInfo.Email)
		switch {
		case errors.Is(err, errx.ErrorUserNotFound):
			ape.RenderErr(w, problems.NotFound("user with this email not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return

	}

	s.log.Infof("User %s logged in with Google", userInfo.Email)

	ape.Render(w, http.StatusOK, responses.TokensPair(tokensPair))
}
