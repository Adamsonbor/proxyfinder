package serviceapiv1

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"proxyfinder/internal/domain"
	"proxyfinder/internal/domain/dto"
	"proxyfinder/pkg/filter"
)

const (
	ErrInvalidField   = "Invalid field name"
	ErrRecordNotFound = "Record not found"
	ErrAlreadyExists  = "Already exists"
	ErrIdNotFound     = "Id not found"
)

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// TODO: fix filter dependency
type ProxyService interface {
	GetAll(ctx context.Context, options filter.Options) ([]dto.Proxy, error)
	Update(ctx context.Context, options filter.Options) error
}

// TODO: fix filter dependency
type FavoritsService interface {
	GetAll(ctx context.Context, options filter.Options) ([]domain.Favorits, error)
}

type JWTService interface {
	JWTMiddleware(next http.Handler) http.Handler
	ParseToken(tokenString string) (*jwt.Token, error)
	GenerateAccessToken(userId int64) (string, error)
	GenerateRefreshToken() (string, error)
	ExtractToken(r *http.Request) (string, error)
	ValidateToken(tokenString string) error
}

type GoogleAuthService interface {
	Login(state string) string
	UpdateRefreshToken(ctx context.Context, refreshToken string) (RefreshResponse, error)
	Callback(ctx context.Context, googleCode string) (domain.User, error)
	GenerateTokens(user domain.User) (string, string, error)
}

// TODO: fix filter dependency
type UserService interface {
	// GetAll(ctx context.Context, options filter.Options) ([]domain.User, error)
	GetBy(ctx context.Context, fieldName string, value interface{}) (domain.User, error)
	Save(ctx context.Context, user domain.User) (int64, error)
	NewSession(ctx context.Context, userId int64, refresh string) error
}

// TODO: fix filter dependency
type ProxyStorage interface {
	GetAll(ctx context.Context, options filter.Options) ([]dto.Proxy, error)
	Update(ctx context.Context, options filter.Options) error
}

// TODO: fix filter dependency
type FavoritsStorage interface {
	GetAll(ctx context.Context, options filter.Options) ([]domain.Favorits, error)
}

// TODO: fix filter dependency
type UserStorage interface {
	GetBy(ctx context.Context, fieldName string, value interface{}) (domain.User, error)
	GetByRefreshToken(ctx context.Context, token string) (domain.User, error)
	Save(ctx context.Context, user domain.User) (int64, error)
	NewSession(ctx context.Context, userId int64, refresh string) error
}
