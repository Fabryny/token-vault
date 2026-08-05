package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Fabryny/token-vault/backend/internal/domain"
)

type Login struct {
	users     domain.UserRepository
	hasher    domain.PasswordHasher
	issuer    domain.TokenIssuer
	dummyHash string
}

func NewLogin(
	users domain.UserRepository,
	hasher domain.PasswordHasher,
	issuer domain.TokenIssuer,
) (*Login, error) {
	// Hash descartável, gerado uma vez no boot. Serve para gastar o MESMO
	// tempo quando o e-mail não existe — ver comentário no Execute.
	dummy, err := hasher.Hash("constant-time-placeholder")
	if err != nil {
		return nil, fmt.Errorf("dummy hash: %w", err)
	}
	return &Login{users: users, hasher: hasher, issuer: issuer, dummyHash: dummy}, nil
}

func (uc *Login) Execute(ctx context.Context, email, password string) (string, error) {
	u, err := uc.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			// Compara contra um hash falso mesmo sem usuário: sem isto, o
			// e-mail inexistente responde MUITO mais rápido (não roda bcrypt),
			// e o tempo de resposta vira enumeração de usuários.
			_ = uc.hasher.Compare(uc.dummyHash, password)
			return "", domain.ErrInvalidCredentials
		}
		return "", fmt.Errorf("find user: %w", err)
	}

	if err := uc.hasher.Compare(u.PasswordHash, password); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := uc.issuer.Issue(u.ID)
	if err != nil {
		return "", fmt.Errorf("issue token: %w", err)
	}
	return token, nil
}
