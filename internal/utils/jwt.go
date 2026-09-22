package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zekielmp/Bitly/internal/config"
)

// Claims contains data for the user
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenTokenPair generates access and refresh token
func GenTokenPair(cfg *config.JwtConfig, userID uint, email, role string) (accessToken, refreshToken string, err error) {
	//Access Token
	accessClaims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := at.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return "", "", err
	}

	//Refresh Token
	refreshClaims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.RefreshTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := rt.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// ValidateToken checks if token is valid
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if Claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return Claims, nil
	}
	return nil, errors.New("Invalid token")
}
