package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTService interface {
	GenerateAccessToken(expirationTime time.Duration, claim JwtCustomClaim) (string, error)
	GenerateRefreshToken(expirationTime time.Duration, claim JwtCustomClaim, jti string) (string, error)
	ValidateToken(token string) (*JwtCustomClaim, error)
}

type JwtCustomClaim struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secret string
}

func NewJWTService(secret string) JWTService {
	return &jwtService{
		secret: secret,
	}
}

func (s *jwtService) GenerateAccessToken(expirationTime time.Duration, claim JwtCustomClaim) (string, error) {
	expiration := time.Now().Add(expirationTime)
	claim.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Now())
	claim.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expiration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claim)
	return token.SignedString([]byte(s.secret))
}

func (s *jwtService) GenerateRefreshToken(expirationTime time.Duration, claim JwtCustomClaim, jti string) (string, error) {
	expiration := time.Now().Add(expirationTime)
	claim.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Now())
	claim.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expiration)
	claim.RegisteredClaims.ID = jti

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claim)
	return token.SignedString([]byte(s.secret))
}

func (s *jwtService) ValidateToken(tokenStr string) (*JwtCustomClaim, error) {
	claims := &JwtCustomClaim{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.secret), nil
	})
	return claims, err
}
