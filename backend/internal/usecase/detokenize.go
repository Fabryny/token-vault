package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Fabryny/token-vault/backend/internal/domain"
	"github.com/google/uuid"
)

// Guarda as ferramentas necessárias para destokenizar (banco e criptografia)
type Detokenize struct {
	repo domain.TokenRepository
	enc  domain.Encryptor
}

// Injeta as dependências necessárias para criar o caso de uso
func NewDetokenize(repo domain.TokenRepository, enc domain.Encryptor) *Detokenize {
	return &Detokenize{repo: repo, enc: enc}
}

// Executa o fluxo completo: busca no banco, valida expiração e descriptografa
func (uc *Detokenize) Execute(ctx context.Context, value string, requesterID uuid.UUID) (string, error) {
	t, err := uc.repo.FindByValue(ctx, value)
	if err != nil {
		return "", err
	}

	// AUTORIZAÇÃO
	if t.OwnerID != requesterID {
		// ErrTokenNotFound de propósito, não "forbidden": responder 403
		// confirmaria que o token EXISTE, só que de outra pessoa.
		return "", domain.ErrTokenNotFound
	}

	if t.IsExpired(time.Now()) {
		return "", domain.ErrTokenExpired
	}

	pan, err := uc.enc.Decrypt(t.Ciphertext, t.Nonce)
	if err != nil {
		return "", fmt.Errorf("decrypt token %s: %w", t.Value, err)
	}
	return string(pan), nil
}
