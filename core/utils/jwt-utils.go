package utils

import (
	core_constants "buckly-ms/core/constants"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId            int64  `json:"user_id"`
	Email             string `json:"email"`
	PhoneNo           string `json:"phone_no"`
	IsPhoneNoVerified bool   `json:"is_phone_no_verified"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId int64, email string, phoneNo string, isPhoneNoVerified bool, secret string) (string, error) {
	claims := Claims{
		UserId:            userId,
		Email:             email,
		PhoneNo:           phoneNo,
		IsPhoneNoVerified: isPhoneNoVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userId int64, secret string) (string, error) {
	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenStr string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

// ValidateTokenIgnoreExpiration validates token claims and signature but allows expired tokens
// Useful for refresh token endpoints where expired access tokens need to be validated
func ValidateTokenIgnoreExpiration(tokenStr string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	// Note: We skip the token.Valid check which would reject expired tokens
	// We only verify the signature was valid by checking if parsing succeeded
	return claims, nil
}

func GetClaims(c *gin.Context) (*Claims, error) {
	claims, exists := c.Get(core_constants.JWT_TOKEN_CLAIM)
	if !exists {
		return nil, errors.New("Claims not found")
	}
	return claims.(*Claims), nil
}
