package jwtservice

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"proxyfinder/internal/api"
	"proxyfinder/internal/config"
	"proxyfinder/internal/storage"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken  = fmt.Errorf("invalid token")
	ErrExpiredToken  = fmt.Errorf("expired token")
	ErrMissingToken  = fmt.Errorf("missing token")
	ErrSigningMethod = fmt.Errorf("invalid signing method")
)

// tokens and session_id
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTService struct {
	secret      string
	log         *slog.Logger
	cfg         *config.Config
	userStorage storage.UserStorage
}

func NewJWTService(
	log *slog.Logger,
	cfg *config.Config,
	userStorage storage.UserStorage,
) *JWTService {
	return &JWTService{
		secret:      cfg.JWT.Secret,
		log:         log,
		cfg:         cfg,
		userStorage: userStorage,
	}
}

// Generate access token
func (self *JWTService) GenerateAccessToken(userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.FormatInt(userId, 10),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(self.cfg.JWT.AccessTokenTTL).Unix(),
	})

	return token.SignedString([]byte(self.secret))
}

// Generate refresh token
func (self *JWTService) GenerateRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uuid.New().String(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(self.cfg.JWT.RefreshTokenTTL).Unix(),
	})

	return token.SignedString([]byte(self.secret))
}

// Extract token from header
func (self *JWTService) ExtractToken(r *http.Request) (string, error) {
	const BEARER_SCHEMA = "Bearer "
	bearToken := r.Header.Get("Authorization")

	return strings.TrimPrefix(bearToken, BEARER_SCHEMA), nil
}

func (self *JWTService) ValidateToken(tokenString string) error {
	token, err := self.ParseToken(tokenString)
	if err != nil {
		return err
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

// String token to jwt.Token
func (self *JWTService) ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrSigningMethod
		}
		return []byte(self.secret), nil
	})
}

// Validate token from header and set user_id int64 in context
func (self *JWTService) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := self.log.With(slog.String("op", "JWTMiddleware"))

		tokenString, err := self.ExtractToken(r)
		if err != nil {
			log.Error("extract", slog.Any("err", err))
			ReturnError(log, w, http.StatusUnauthorized, ErrMissingToken)
			return
		}
		log.Debug("token string", slog.Any("token", tokenString))

		token, err := self.ParseToken(tokenString)
		if err != nil {
			log.Error("parse", slog.Any("err", err))
			ReturnError(log, w, http.StatusUnauthorized, ErrInvalidToken)
			return
		}
		log.Debug("token", slog.Any("token", token))

		if !token.Valid {
			log.Error("invalid")
			ReturnError(log, w, http.StatusUnauthorized, ErrInvalidToken)
			return
		}

		sUserId, err := token.Claims.GetSubject()
		if err != nil {
			log.Error("claims", slog.Any("err", err))
			ReturnError(log, w, http.StatusUnauthorized, ErrInvalidToken)
			return
		}
		log.Debug("user_id", slog.Any("user_id", sUserId))

		i64UserId, err := strconv.ParseInt(sUserId, 10, 64)
		if err != nil {
			log.Error("Atoi", slog.Any("err", err))
			ReturnError(log, w, http.StatusInternalServerError, err)
			return
		}

		if i64UserId == 0 {
			log.Error("user_id is zero")
			ReturnError(log, w, http.StatusUnauthorized, ErrInvalidToken)
			return
		}

		// set user_id in context
		ctx := context.WithValue(r.Context(), "user_id", i64UserId)
		log.Info("success", slog.String("user_id", sUserId))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ReturnError(log *slog.Logger, w http.ResponseWriter, statusCode int, err error) {
	log.Error("error", slog.Any("error", err))
	w.WriteHeader(statusCode)
	api.ReturnResponse(w, "error", nil, err.Error())
}
