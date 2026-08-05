package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Fabryny/token-vault/backend/internal/domain"
	"github.com/Fabryny/token-vault/backend/internal/usecase"
)

type TokenHandler struct {
	tokenize   *usecase.Tokenize
	detokenize *usecase.Detokenize
	log        *zap.Logger
}

func NewTokenHandler(t *usecase.Tokenize, d *usecase.Detokenize, log *zap.Logger) *TokenHandler {
	return &TokenHandler{tokenize: t, detokenize: d, log: log}
}

type tokenizeRequest struct {
	// PAN real tem 13 a 19 dígitos — não fixar em 16.
	// "number" = ^[0-9]+$ (só dígitos). NÃO usar "numeric",
	// que aceita sinal e decimal: ^[-+]?[0-9]+(?:\.[0-9]+)?$
	PAN string `json:"pan" binding:"required,number,min=13,max=19"`
}

type tokenizeResponse struct {
	Token string `json:"token"`
	Last4 string `json:"last4"`
}

func (h *TokenHandler) getOwnerID(c *gin.Context) (uuid.UUID, bool) {
	ownerID, ok := userIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}
	return ownerID, true
}

func (h *TokenHandler) Tokenize(c *gin.Context) {
	ownerID, ok := h.getOwnerID(c)
	if !ok {
		return
	}

	var req tokenizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// NUNCA err.Error() aqui: a mensagem do validator pode ecoar o PAN.
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	t, err := h.tokenize.Execute(c.Request.Context(), req.PAN, ownerID)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPAN) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card number"})
			return
		}
		h.log.Error("tokenize failed", zap.Error(err)) // sem PAN
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	h.log.Info("token created",
		zap.String("token", t.Value),
		zap.String("owner_id", t.OwnerID.String()),
	)

	c.JSON(http.StatusCreated, tokenizeResponse{Token: t.Value, Last4: t.Last4})
}

type detokenizeRequest struct {
	Token string `json:"token" binding:"required"`
}

type detokenizeResponse struct {
	PAN string `json:"pan"`
}

func (h *TokenHandler) Detokenize(c *gin.Context) {
	var req detokenizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ownerID, ok := h.getOwnerID(c)
	if !ok {
		return
	}

	pan, err := h.detokenize.Execute(c.Request.Context(), req.Token, ownerID)
	if err != nil {
		// Mesma resposta para "não existe" e "expirou": a diferença é
		// informação útil para quem está sondando tokens.
		if errors.Is(err, domain.ErrTokenNotFound) || errors.Is(err, domain.ErrTokenExpired) {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		h.log.Error("detokenize failed", zap.String("token", req.Token), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, detokenizeResponse{PAN: pan})
}
