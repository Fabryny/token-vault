package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Fabryny/token-vault/backend/internal/domain"
)

type Register struct {
	users  domain.UserRepository
	hasher domain.PasswordHasher
}

func NewRegister(users domain.UserRepository, hasher domain.PasswordHasher) *Register {
	return &Register{users: users, hasher: hasher}
}

func (uc *Register) Execute(ctx context.Context, email, password string) error {
	hash, err := uc.hasher.Hash(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	u := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	if err := uc.users.Save(ctx, u); err != nil {
		return err // ErrEmailTaken já vem traduzido
	}
	return nil
}
