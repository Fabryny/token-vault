package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
)

// Guarda o motor de criptografia pronto na memória para não ter que recriar toda hora
type Encryptor struct {
	aead cipher.AEAD
}

// Prepara e valida o motor usando a sua chave secreta do .env
func NewEncryptor(key []byte) (*Encryptor, error) {
	// Valida o tamanho da chave e cria a base do AES
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	// Configura o padrão GCM que embaralha e também impede que alguém altere o dado no banco
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	return &Encryptor{aead: aead}, nil
}

// Recebe texto limpo e devolve texto blindado e a senha temporária usada
func (e *Encryptor) Encrypt(plaintext []byte) (ciphertext, nonce []byte, err error) {
	// Cria o nonce, uma pequena senha aleatória obrigatória que deve ser única por dado
	nonce = make([]byte, e.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("generate nonce: %w", err)
	}

	// Seal faz o trabalho de trancar o texto usando a sua chave e o nonce
	return e.aead.Seal(nil, nonce, plaintext, nil), nonce, nil
}

// Recebe o texto blindado e o nonce para devolver o texto limpo original
func (e *Encryptor) Decrypt(ciphertext, nonce []byte) ([]byte, error) {
	// Open tenta destrancar. Se o dado no banco foi modificado ou a chave errada, ele falha
	plaintext, err := e.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// Oculta o motivo real do erro de propósito para não dar pistas a hackers
		return nil, errors.New("decrypt failed")
	}

	return plaintext, nil
}
