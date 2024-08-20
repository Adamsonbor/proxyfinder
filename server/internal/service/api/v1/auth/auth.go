package auth

import (
	"net/http"
)

type OAuth2 interface {
	Middleware(next http.Handler) http.Handler
	Login(w http.ResponseWriter, r *http.Request)
	Callback(w http.ResponseWriter, r *http.Request)
	UserInfo(w http.ResponseWriter, r *http.Request)
}

type JWTService interface {
	JWTMiddleware(next http.Handler) http.Handler
	GenerateAccessToken(userId int64) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateToken(tokenString string) error
	ExtractToken(r *http.Request) (string, error)
}
