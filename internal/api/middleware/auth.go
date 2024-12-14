package middleware

import (
	"errors"
	"fmt"
	"time"

	"auth/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

var (
	secretKey       = config.SecretKey
	accessTokenexp  = time.Duration(config.AccessTokenexp)
	refreshTokenexp = time.Duration(config.RefreshTokenexp)

	// Custom errors for more precise error handling
	ErrTokenExpired = errors.New("token has expired")
	ErrInvalidToken = errors.New("invalid token")
	ErrMissingToken = errors.New("authorization token is missing")
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID    int       //`json:"user_id"`
	TokenType TokenType //`json:"token_type"`
}

func GenerateToken(userID int, tokenType TokenType, duration time.Duration) (string, error) {
	// Validate input
	if userID <= 0 {
		return "", fmt.Errorf("invalid user ID")
	}

	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "",
		},
		UserID:    userID,
		TokenType: tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Error().Err(err).Int("userID", userID).Str("tokenType", string(tokenType)).
			Msg("Failed to generate token")
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

func GenerateAccessToken(userID int) (string, error) {
	return GenerateToken(userID, AccessToken, accessTokenexp)
}

func GenerateRefreshToken(userID int) (string, error) {
	return GenerateToken(userID, RefreshToken, refreshTokenexp)
}

func ValidateToken(tokenString string) (*TokenClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		log.Error().Err(err).Msg("Token validation failed")
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		log.Warn().Msg("Invalid token claims")
		return nil, ErrInvalidToken
	}

	return claims, nil
}
