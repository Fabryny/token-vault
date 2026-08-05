package domain

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Token é o registro do vault: o mapeamento token → PAN cifrado.
type Token struct {
	ID         uuid.UUID
	Value      string // o token aleatório — o que sai para fora
	Ciphertext []byte // o PAN cifrado (AES-GCM)
	Nonce      []byte // nonce do AES-GCM, ÚNICO por registro
	Last4      string // 4 últimos dígitos, exibíveis sem detokenizar
	OwnerID    uuid.UUID
	CreatedAt  time.Time
	ExpiresAt  *time.Time // nil = não expira
}

// Erros de domínio. A infra TRADUZ os erros dela para estes,
// para que o usecase nunca precise importar pgx.
var (
	ErrInvalidPAN    = errors.New("invalid PAN")
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token expired")
)

func NewTokenValue() string {
	return rand.Text()
}

func (t *Token) IsExpired(now time.Time) bool {
	return t.ExpiresAt != nil && now.After(*t.ExpiresAt)
}
