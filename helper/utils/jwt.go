package utils

import (
	"arek-muhammadiyah-be/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID string `json:"user_id"`
	RoleID *uint  `json:"role_id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken - Generate access token dengan expiry 1 jam
func GenerateAccessToken(userID string, roleID *uint) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// GenerateRefreshToken - Generate refresh token dengan expiry 30 hari
func GenerateRefreshToken(userID string) (string, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour) // 30 hari

	claims := &RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTRefreshSecret))
}

// ValidateAccessToken - Validasi access token
func ValidateAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

// ValidateRefreshToken - Validasi refresh token
func ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTRefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}

// Legacy function untuk backward compatibility - DEPRECATED, gunakan ValidateAccessToken
func ValidateToken(tokenString string) (*Claims, error) {
	return ValidateAccessToken(tokenString)
}

// Legacy function untuk backward compatibility - DEPRECATED, gunakan GenerateAccessToken
func GenerateToken(userID string, roleID *uint) (string, error) {
	return GenerateAccessToken(userID, roleID)
}