package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManagerInterface defines the interface for JWT operations
type JWTManagerInterface interface {
	GenerateToken(userID, email string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	ValidateRefreshToken(tokenString string) (*Claims, error)
}

// JWTManager handles JWT token operations
type JWTManager struct {
	secretKey     string
	refreshSecret string
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
}

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey, refreshSecret string, tokenExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		refreshSecret: refreshSecret,
		tokenExpiry:   tokenExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateToken generates a new JWT token
func (j *JWTManager) GenerateToken(userID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// GenerateRefreshToken generates a new refresh token
func (j *JWTManager) GenerateRefreshToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.refreshSecret))
}

// ValidateToken validates a JWT token
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateRefreshToken validates a refresh token
func (j *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.refreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}
