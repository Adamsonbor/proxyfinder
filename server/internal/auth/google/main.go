package googleauth

import (
	"net/http"
	"proxyfinder/internal/config"
	"time"

	"golang.org/x/oauth2"
)

type GoogleAuth struct {
	cfg *oauth2.Config
	env *config.GoogleAuth
}

func NewGoogleAuth(env *config.GoogleAuth) *GoogleAuth {
	return &GoogleAuth{
		env: env,
		cfg: &oauth2.Config{
			ClientID:     env.ClientID,
			ClientSecret: env.ClientSecret,
			Scopes:       env.Scope,
			Endpoint: oauth2.Endpoint{
				AuthURL:  env.Endpoint.AuthURL,
				TokenURL: env.Endpoint.TokenURL,
			},
			RedirectURL: env.RedirectURL,
		},
	}
}

func (g *GoogleAuth) Login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, g.cfg.AuthCodeURL("state"), http.StatusFound)
}

func (g *GoogleAuth) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (g *GoogleAuth) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, err := g.cfg.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token.AccessToken,
		Path:  "/",
		MaxAge: int(token.Expiry.Sub(time.Now()).Seconds()),
		Expires: token.Expiry,
	})
	http.Redirect(w, r, g.env.HomeUrl, http.StatusFound)
}

func (g *GoogleAuth) TokenOrLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("token")
		if err != nil {
			g.Login(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
