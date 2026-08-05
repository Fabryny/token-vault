package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Fabryny/token-vault/backend/internal/domain"
)

type TokenRepository struct {
	pool *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{pool: pool}
}

const saveTokenSQL = `
INSERT INTO tokens (id, token, ciphertext, nonce, last4, owner_id, created_at, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

func (r *TokenRepository) Save(ctx context.Context, t *domain.Token) error {
	_, err := r.pool.Exec(ctx, saveTokenSQL,
		t.ID, t.Value, t.Ciphertext, t.Nonce, t.Last4, t.OwnerID, t.CreatedAt, t.ExpiresAt)
	if err != nil {
		// %w com o err do pgx: ele NÃO contém os valores dos parâmetros.
		return fmt.Errorf("insert token: %w", err)
	}
	return nil
}

const findByValueSQL = `
SELECT id, token, ciphertext, nonce, last4, owner_id, created_at, expires_at
FROM tokens
WHERE token = $1`

func (r *TokenRepository) FindByValue(ctx context.Context, value string) (*domain.Token, error) {
	var t domain.Token

	err := r.pool.QueryRow(ctx, findByValueSQL, value).Scan(
		&t.ID, &t.Value, &t.Ciphertext, &t.Nonce,
		&t.Last4, &t.OwnerID, &t.CreatedAt, &t.ExpiresAt,
	)

	// A TRADUÇÃO acontece aqui, na borda: erro de infra vira erro de domínio.
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrTokenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query token: %w", err)
	}

	return &t, nil
}
