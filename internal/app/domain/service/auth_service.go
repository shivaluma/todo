package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/domain/repository"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// AuthService provides authentication related functionality
type AuthService struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
	jwtKey   []byte
	jwtExp   time.Duration
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Fullname string `json:"fullname"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repository.UserRepository, logger *logger.Logger, jwtSecret string, jwtExpiration time.Duration) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		logger:   logger,
		jwtKey:   []byte(jwtSecret),
		jwtExp:   jwtExpiration,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, fullname, email, password string) (*model.User, error) {
	// Check if user already exists
	exists, err := s.userRepo.Exists(ctx, email)
	if err != nil {
		s.logger.Error("Failed to check if user exists", "error", err)
		return nil, err
	}

	if exists {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	user, err := model.NewUser(fullname, email, password)
	if err != nil {
		s.logger.Error("Failed to create user", "error", err)
		return nil, err
	}

	// Save user to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to save user", "error", err)
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("User not found", "email", email)
		return "", errors.New("invalid credentials")
	}

	// Check password
	if !user.CheckPassword(password) {
		s.logger.Error("Invalid password", "email", email)
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		s.logger.Error("Failed to generate token", "error", err)
		return "", err
	}

	return token, nil
}

// GenerateToken generates a JWT token for a user
func (s *AuthService) GenerateToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(s.jwtExp)

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
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetUserFromToken gets a user from a JWT token
func (s *AuthService) GetUserFromToken(ctx context.Context, tokenString string) (*model.User, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) GetUserFromId(ctx context.Context, id string) (*model.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
