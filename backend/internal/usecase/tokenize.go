package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Fabryny/token-vault/backend/internal/domain"
)

type Tokenize struct {
	repo domain.TokenRepository
	enc  domain.Encryptor
}

func NewTokenize(repo domain.TokenRepository, enc domain.Encryptor) *Tokenize {
	return &Tokenize{repo: repo, enc: enc}
}

// Execute valida → gera token → criptografa → salva. Devolve SÓ o token.
func (uc *Tokenize) Execute(ctx context.Context, pan string, ownerID uuid.UUID) (*domain.Token, error) {
	if !domain.IsValidLuhn(pan) {
		return nil, domain.ErrInvalidPAN
	}

	ciphertext, nonce, err := uc.enc.Encrypt([]byte(pan))
	if err != nil {
		return nil, fmt.Errorf("encrypt pan: %w", err)
	}

	t := &domain.Token{
		ID:         uuid.New(),
		Value:      domain.NewTokenValue(),
		Ciphertext: ciphertext,
		Nonce:      nonce,
		Last4:      pan[len(pan)-4:],
		OwnerID:    ownerID,
		CreatedAt:  time.Now(),
	}

	if err := uc.repo.Save(ctx, t); err != nil {
		return nil, fmt.Errorf("save token: %w", err)
	}

	return t, nil
}
