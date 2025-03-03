package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Fullname string `json:"fullname"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	secretKey []byte
	expiration time.Duration
	logger    *logger.Logger
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string, tokenDuration time.Duration, logger *logger.Logger) *JWTService {
	return &JWTService{
		secretKey:  []byte(secretKey),
		expiration: tokenDuration,
		logger:     logger,
	}
}

// GenerateToken generates a JWT token for a user
func (s *JWTService) GenerateToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(s.expiration)

	claims := &JWTClaims{
		UserID:   user.ID.String(),
		Fullname: user.Fullname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "todo-api",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		s.logger.Error("Failed to sign JWT token", "error", err)
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		s.logger.Error("Failed to parse JWT token", "error", err)
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractUserID extracts the user ID from a JWT token
func (s *JWTService) ExtractUserID(tokenString string) (uuid.UUID, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		s.logger.Error("Failed to parse user ID from JWT token", "error", err)
		return uuid.Nil, err
	}

	return userID, nil
}

// RefreshToken refreshes a JWT token
func (s *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Create a new token with the same claims but a new expiration time
	expirationTime := time.Now().Add(s.expiration)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.NotBefore = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		s.logger.Error("Failed to sign refreshed JWT token", "error", err)
		return "", err
	}

	return newTokenString, nil
}
