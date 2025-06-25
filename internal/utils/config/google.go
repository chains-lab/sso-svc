package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func (c *Config) GoogleOAuth() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.OAuth.Google.ClientID,
		ClientSecret: c.OAuth.Google.ClientSecret,
		RedirectURL:  c.OAuth.Google.RedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}
