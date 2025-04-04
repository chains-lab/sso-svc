package handlers

import (
	"net/http"

	"golang.org/x/oauth2"
)

func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.google.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
