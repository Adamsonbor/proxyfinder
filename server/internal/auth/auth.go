package auth

import "net/http"

type Auth interface {
	AuthMiddleware(next http.Handler) http.Handler
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Callback(w http.ResponseWriter, r *http.Request)
}
