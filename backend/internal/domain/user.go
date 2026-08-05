package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string // bcrypt
	CreatedAt    time.Time
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already registered")
	ErrUnauthorized       = errors.New("unauthorized")
)
