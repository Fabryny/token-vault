CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,          -- bcrypt, NUNCA a senha
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE tokens (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    token       TEXT        NOT NULL UNIQUE,     -- o token aleatório (o que sai para fora)
    ciphertext  BYTEA       NOT NULL,            -- o PAN criptografado (AES-GCM)
    nonce       BYTEA       NOT NULL,            -- nonce do AES-GCM, ÚNICO por registro
    last4       CHAR(4)     NOT NULL,            -- últimos 4 dígitos, exibir sem detokenizar
    owner_id    UUID        NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at  TIMESTAMPTZ
);

CREATE INDEX idx_tokens_owner_created ON tokens (owner_id, created_at DESC);

-- Usuário de DESENVOLVIMENTO.
INSERT INTO users (id, email, password_hash) VALUES
    ('00000000-0000-0000-0000-000000000001', 'dev@local', 'NAO_E_UM_HASH_VALIDO');