package handlers

import (
	"net/http"

	"golang.org/x/oauth2"
)

func (a *App) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := a.GoogleOAuth.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
