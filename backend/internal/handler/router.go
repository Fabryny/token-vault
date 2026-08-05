package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/Fabryny/token-vault/backend/internal/domain"
)

func RegisterRoutes(r *gin.Engine, t *TokenHandler, a *AuthHandler, issuer domain.TokenIssuer) {
	api := r.Group("/api")

	// Públicas
	api.POST("/auth/register", a.Register)
	api.POST("/auth/login", a.Login)
	api.POST("/auth/logout", a.Logout) // apagar cookie não exige sessão válida

	// Protegidas
	protected := api.Group("", AuthMiddleware(issuer))
	protected.GET("/auth/me", a.Me) // o front pergunta "minha sessão vale?"
	protected.POST("/tokenize", t.Tokenize)
	protected.POST("/detokenize", t.Detokenize)
}
