package domain

import (
	"context"

	"github.com/google/uuid"
)

// TokenRepository é implementado pela infra (Postgres).
// Declarado AQUI, no domínio, porque é o domínio quem consome.
type TokenRepository interface {
	Save(ctx context.Context, t *Token) error
	FindByValue(ctx context.Context, value string) (*Token, error)
}

// Encryptor é implementado por internal/infra/crypto.
type Encryptor interface {
	Encrypt(plaintext []byte) (ciphertext, nonce []byte, err error)
	Decrypt(ciphertext, nonce []byte) ([]byte, error)
}

type UserRepository interface {
	Save(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// PasswordHasher é implementado por infra/crypto (bcrypt).
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

// TokenIssuer é implementado por infra/auth (JWT).
type TokenIssuer interface {
	Issue(userID uuid.UUID) (string, error)
	Parse(raw string) (uuid.UUID, error)
}
