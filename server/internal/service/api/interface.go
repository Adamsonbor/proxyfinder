package serviceapiv1

import (
	"context"
	"net/http"
	"proxyfinder/internal/domain"
	"proxyfinder/internal/domain/dto"
	"proxyfinder/pkg/options"

	"github.com/golang-jwt/jwt/v5"
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
	ExpiresInRef int64  `json:"expires_in_ref"`
}

type JWTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	ExpiresInRef int64  `json:"expires_in_ref"`
}

// TODO: fix filter dependency
type ProxyService interface {
	GetAll(ctx context.Context, filter options.Options, sort options.Options) ([]dto.Proxy, error)
	Update(ctx context.Context, filter options.Options) error
}

// TODO: fix filter dependency
type FavoritsService interface {
	GetAll(ctx context.Context, filter options.Options, sort options.Options) ([]domain.Favorits, error)
	Save(ctx context.Context, options options.Options) (int64, error)
	Delete(ctx context.Context, options options.Options) error
}

type JWTService interface {
	// get access_toke and set user_id in context
	JWTMiddleware(next http.Handler) http.Handler
	ParseToken(tokenString string) (*jwt.Token, error)
	GenerateAccessToken(userId int64) (*jwt.Token, error)
	GenerateRefreshToken() (*jwt.Token, error)
	ExtractToken(r *http.Request) (string, error)
	ValidateToken(tokenString string) error
}

type GoogleAuthService interface {
	Login(state string) string
	UpdateRefreshToken(ctx context.Context, refreshToken string) (*JWTokens, error)
	Callback(ctx context.Context, googleCode string) (*JWTokens, error)
	GenerateTokens(user domain.User) (*jwt.Token, *jwt.Token, error)
}

// TODO: fix filter dependency
type UserService interface {
	// GetAll(ctx context.Context, options filter.Options) ([]domain.User, error)
	UserInfo(ctx context.Context, id int64) (domain.User, error)
	GetBy(ctx context.Context, fieldName string, value interface{}) (domain.User, error)
	Save(ctx context.Context, user domain.User) (int64, error)
	NewSession(ctx context.Context, userId int64, token string, expiresIn int64) error
}

// TODO: fix filter dependency
type ProxyStorage interface {
	GetAll(ctx context.Context, filter options.Options, sort options.Options) ([]dto.Proxy, error)
	Update(ctx context.Context, filter options.Options) error
}

// TODO: fix filter dependency
type FavoritsStorage interface {
	GetAll(ctx context.Context, filter options.Options, sort options.Options) ([]domain.Favorits, error)
	Save(ctx context.Context, options options.Options) (int64, error)
	Delete(ctx context.Context, options options.Options) error
}

// TODO: fix filter dependency
type UserStorage interface {
	GetBy(ctx context.Context, fieldName string, value interface{}) (domain.User, error)
	GetByRefreshToken(ctx context.Context, token string) (domain.User, error)
	Save(ctx context.Context, user domain.User) (int64, error)
	NewSession(ctx context.Context, userId int64, refresh string, expiresAt int64) error
}
