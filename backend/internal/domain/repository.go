package domain

import "context"

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
