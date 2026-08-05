package crypto

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// maxPasswordBytes: bcrypt trunca/rejeita acima de 72 BYTES.
const maxPasswordBytes = 72

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{cost: bcrypt.DefaultCost} // 10 (Min 4, Max 31)
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	if len(password) > maxPasswordBytes {
		return "", fmt.Errorf("senha acima de %d bytes", maxPasswordBytes)
	}
	b, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(b), nil
}

func (h *BcryptHasher) Compare(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errors.New("password mismatch")
	}
	return nil
}
