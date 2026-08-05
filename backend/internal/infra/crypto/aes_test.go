package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Fabryny/token-vault/backend/internal/infra/crypto"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	enc, err := crypto.NewEncryptor(make([]byte, 32))
	require.NoError(t, err)

	pan := "4111111111111111"

	ciphertext, nonce, err := enc.Encrypt([]byte(pan))
	require.NoError(t, err)
	assert.NotEqual(t, pan, string(ciphertext)) // não é o PAN em claro
	assert.Len(t, nonce, 12)

	got, err := enc.Decrypt(ciphertext, nonce)
	require.NoError(t, err)
	assert.Equal(t, pan, string(got))
}

func TestDecryptRejectsTamperedCiphertext(t *testing.T) {
	enc, err := crypto.NewEncryptor(make([]byte, 32))
	require.NoError(t, err)

	ciphertext, nonce, err := enc.Encrypt([]byte("4111111111111111"))
	require.NoError(t, err)

	ciphertext[0] ^= 0xFF // vira um bit

	_, err = enc.Decrypt(ciphertext, nonce)
	assert.Error(t, err) // é isso que o "A" de AEAD garante
}

func TestNonceIsNeverReused(t *testing.T) {
	enc, err := crypto.NewEncryptor(make([]byte, 32))
	require.NoError(t, err)

	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		_, nonce, err := enc.Encrypt([]byte("4111111111111111"))
		require.NoError(t, err)
		assert.False(t, seen[string(nonce)], "nonce repetido!")
		seen[string(nonce)] = true
	}
}
