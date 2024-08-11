package main

import (
	"context"
	"log"
	"net/http"
	"proxyfinder/internal/domain"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

var gCfg = oauth2.Config{
	ClientID:     "592292669736-qrqbci7vqkaqd2rhalot5bsvnb7ek5sl.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes: []string{
		people.UserinfoProfileScope,
		people.UserinfoEmailScope,
	},
	Endpoint: google.Endpoint,
}

func main() {
	log.Printf("Starting server")

	r := chi.NewRouter()

	r.Get("/auth/google/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, gCfg.AuthCodeURL("state"), http.StatusFound)
	})
	r.Get("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		token, err := gCfg.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := getUserInfo(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println(user)
	})

	log.Println("Listening on :8080")
	chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("Route %s %s\n", method, route)
		return nil
	})
	http.ListenAndServe(":8080", r)

}

func getUserInfo(token *oauth2.Token) (*domain.User, error) {
	client := gCfg.Client(context.Background(), token)

	svc, err := people.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	userInfo, err := svc.People.Get("people/me").PersonFields("names,emailAddresses,photos,phoneNumbers").Do()
	if err != nil {
		return nil, err
	}

	log.Println(userInfo)

	return &domain.User{
		Email: userInfo.EmailAddresses[0].Value,
		Name:  userInfo.Names[0].DisplayName,
	}, nil
}
