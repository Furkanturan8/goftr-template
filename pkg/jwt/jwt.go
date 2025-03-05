package jwt

import (
	"goftr-v1/config"
	"goftr-v1/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtConfig *config.JWTConfig

func Init(cfg *config.JWTConfig) {
	jwtConfig = cfg
}

type Claims struct {
	UserID int64      `json:"user_id"`
	Role   model.Role `json:"role"`
	Email  string     `json:"email"`
	jwt.RegisteredClaims
}

func Generate(user *model.User) (string, error) {
	claims := Claims{
		user.ID,
		user.Role,
		user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtConfig.Expiration) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.Secret))
}

func Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
