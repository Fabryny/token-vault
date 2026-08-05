package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Fabryny/token-vault/backend/internal/domain"
)

const ctxUserID = "userID"

func AuthMiddleware(issuer domain.TokenIssuer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// O navegador manda o cookie sozinho. Nada de header Authorization:
		// o token é HttpOnly, então o JavaScript nem sabe que ele existe.
		raw, err := c.Cookie(SessionCookieName)
		if err != nil || raw == "" {
			// AbortWithStatusJSON escreve a resposta E interrompe a cadeia.
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID, err := issuer.Parse(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set(ctxUserID, userID)
		c.Next() // segue para o handler
	}
}

// userIDFrom lê o que o middleware guardou. Sem MustGet: pânico em
// requisição é jeito ruim de descobrir que a rota ficou fora do grupo.
func userIDFrom(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(ctxUserID)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}
