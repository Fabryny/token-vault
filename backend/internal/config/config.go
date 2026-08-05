package config

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

// Config é tudo que a aplicação precisa do ambiente.
type Config struct {
	DatabaseURL   string
	EncryptionKey []byte
	JWTSecret     string
	CookieSecure  bool
	Port          string
}

// Load lê e VALIDA o ambiente. Erro aqui derruba a app no boot, de propósito.
func Load() (*Config, error) {

	jwtSecret := os.Getenv("APP_JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("APP_JWT_SECRET não definida")
	}

	dbURL := os.Getenv("APP_DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("APP_DATABASE_URL não definida")
	}

	rawKey := os.Getenv("APP_ENCRYPTION_KEY")
	if rawKey == "" {
		return nil, errors.New("APP_ENCRYPTION_KEY não definida")
	}

	key, err := hex.DecodeString(rawKey)
	if err != nil {
		return nil, errors.New("APP_ENCRYPTION_KEY não é hex válido")
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("APP_ENCRYPTION_KEY: esperado 32 bytes (64 hex), veio %d", len(key))
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // default só para o que é seguro ter default
	}

	// Secure=true faz o navegador só mandar o cookie por HTTPS.
	// Em dev (http://localhost) isso sumiria com o cookie, então vem do ambiente.
	cookieSecure := os.Getenv("APP_COOKIE_SECURE") == "true"

	return &Config{
		DatabaseURL:   dbURL,
		EncryptionKey: key,
		JWTSecret:     jwtSecret,
		CookieSecure:  cookieSecure,
		Port:          port,
	}, nil
}
