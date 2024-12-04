package common

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("uQ8UWgsfKkCdFXRF4VcmhyTxY9Q2htdk")

type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// GenerateAccessToken generates a new JWT access token
func GenerateAccessToken(userID string) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(), // Refresh token valid for 30 days
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// GenerateRefreshToken generates a new JWT refresh token
func GenerateRefreshToken(userID string) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(90 * 24 * time.Hour).Unix(), // Refresh token valid for 90 days
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken validates a JWT token and returns the claims if valid
func ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("⚠️ invalid token: %v", err)
	}
}
