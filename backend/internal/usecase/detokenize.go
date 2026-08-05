package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Fabryny/token-vault/backend/internal/domain"
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
func (uc *Detokenize) Execute(ctx context.Context, value string) (string, error) {
	// Busca o token salvo no banco de dados pelo seu identificador
	t, err := uc.repo.FindByValue(ctx, value)
	if err != nil {
		return "", err // Retorna o erro original caso não encontre
	}

	// Confere se o token já expirou comparando com a hora atual
	if t.IsExpired(time.Now()) {
		return "", domain.ErrTokenExpired
	}

	// Descriptografa o texto blindado usando o ciphertext e o nonce salvos
	pan, err := uc.enc.Decrypt(t.Ciphertext, t.Nonce)
	if err != nil {
		return "", fmt.Errorf("decrypt token %s: %w", t.Value, err) // O token pode aparecer no log de erro, mas o dado sensível (PAN) nunca
	}

	// Converte os bytes descriptografados de volta para texto e retorna o resultado limpo
	return string(pan), nil
}
