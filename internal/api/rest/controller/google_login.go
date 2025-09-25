package controller

import (
	"net/http"

	"golang.org/x/oauth2"
)

func (s Service) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := s.google.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
