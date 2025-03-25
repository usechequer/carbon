package utilities

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

var store *sessions.CookieStore

func RegisterOauthProviders() {
	goth.UseProviders(google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("GOOGLE_CLIENT_REDIRECT")))
	goth.UseProviders(github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), os.Getenv("GITHUB_CLIENT_REDIRECT")))
}

func GetOauthSessionStore() *sessions.CookieStore {
	if store != nil {
		return store
	}

	store = sessions.NewCookieStore([]byte(os.Getenv("JWT_SECRET")))
	return store
}
