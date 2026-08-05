package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *TokenHandler) {
	api := r.Group("/api")
	api.POST("/tokenize", h.Tokenize)
	api.POST("/detokenize", h.Detokenize)
}
