package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/Fabryny/token-vault/backend/internal/config"
	"github.com/Fabryny/token-vault/backend/internal/handler"
	"github.com/Fabryny/token-vault/backend/internal/infra/auth"
	"github.com/Fabryny/token-vault/backend/internal/infra/crypto"
	"github.com/Fabryny/token-vault/backend/internal/infra/database"
	"github.com/Fabryny/token-vault/backend/internal/usecase"
)

func main() {
	loadDotEnv()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("logger: %v", err)
	}
	defer logger.Sync()

	ctx := context.Background()

	pool, err := database.NewPool(ctx, cfg)
	if err != nil {
		logger.Fatal("database", zap.Error(err))
	}
	defer pool.Close()

	enc, err := crypto.NewEncryptor(cfg.EncryptionKey)
	if err != nil {
		logger.Fatal("encryptor", zap.Error(err))
	}

	// Um só lugar decide o prazo: o JWT e o cookie expiram juntos.
	// Prazos diferentes dariam "cookie válido com token vencido", ou o contrário.
	const sessionTTL = time.Hour

	hasher := crypto.NewBcryptHasher()
	issuer := auth.NewJWTIssuer(cfg.JWTSecret, sessionTTL)

	// Composition root: é AQUI, e só aqui, que as implementações concretas
	// encontram as interfaces. Nenhuma outra camada sabe que existe Postgres.
	tokenRepo := database.NewTokenRepository(pool)
	userRepo := database.NewUserRepository(pool)

	// NewLogin devolve erro porque gera o hash descartável no boot.
	loginUC, err := usecase.NewLogin(userRepo, hasher, issuer)
	if err != nil {
		logger.Fatal("login usecase", zap.Error(err))
	}

	tokenHandler := handler.NewTokenHandler(
		usecase.NewTokenize(tokenRepo, enc),
		usecase.NewDetokenize(tokenRepo, enc),
		logger,
	)

	authHandler := handler.NewAuthHandler(
		usecase.NewRegister(userRepo, hasher),
		loginUC,
		logger,
		cfg.CookieSecure,
		sessionTTL,
	)

	r := gin.Default() // Logger + Recovery embutidos
	handler.RegisterRoutes(r, tokenHandler, authHandler, issuer)

	logger.Info("server starting", zap.String("port", cfg.Port))
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatal("server", zap.Error(err))
	}
}

// loadDotEnv é conveniência de DESENVOLVIMENTO apenas.
func loadDotEnv() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
}
