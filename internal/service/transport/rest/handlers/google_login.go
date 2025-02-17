package handlers

import (
	"net/http"

	"golang.org/x/oauth2"
)

func (h *Handlers) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	svc := h.svc

	url := svc.GoogleOAuth.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
