package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUser creates a new user
func NewUser(username, email, password string) (*User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// CheckPassword checks if the provided password matches the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// UpdatePassword updates the user's password
func (u *User) UpdatePassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hashedPassword)
	u.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateEmail updates the user's email
func (u *User) UpdateEmail(email string) {
	u.Email = email
	u.UpdatedAt = time.Now().UTC()
}

// UpdateUsername updates the user's username
func (u *User) UpdateUsername(username string) {
	u.Username = username
	u.UpdatedAt = time.Now().UTC()
}
