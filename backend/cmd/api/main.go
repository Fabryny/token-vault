package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/Fabryny/token-vault/backend/internal/config"
	"github.com/Fabryny/token-vault/backend/internal/handler"
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

	// Composition root: é AQUI, e só aqui, que as implementações concretas
	// encontram as interfaces. Nenhuma outra camada sabe que existe Postgres.
	repo := database.NewTokenRepository(pool)
	h := handler.NewTokenHandler(
		usecase.NewTokenize(repo, enc),
		usecase.NewDetokenize(repo, enc),
		logger,
	)

	r := gin.Default() // Logger + Recovery embutidos
	handler.RegisterRoutes(r, h)

	logger.Info("server starting", zap.String("port", cfg.Port))
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatal("server", zap.Error(err))
	}
}

// loadDotEnv é conveniência de DESENVOLVIMENTO apenas.
func loadDotEnv() {
	_ = godotenv.Load() // rodando da raiz, caso ache ignora a próxima.
	_ = godotenv.Load("../.env")
}
