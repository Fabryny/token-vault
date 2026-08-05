package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Fabryny/token-vault/backend/internal/domain"
	"github.com/Fabryny/token-vault/backend/internal/usecase"
)

// SessionCookieName é usado pelo handler (escreve) e pelo middleware (lê).
const SessionCookieName = "tv_session"

type AuthHandler struct {
	register     *usecase.Register
	login        *usecase.Login
	log          *zap.Logger
	cookieSecure bool
	sessionTTL   time.Duration
}

func NewAuthHandler(
	r *usecase.Register,
	l *usecase.Login,
	log *zap.Logger,
	cookieSecure bool,
	sessionTTL time.Duration,
) *AuthHandler {
	return &AuthHandler{
		register:     r,
		login:        l,
		log:          log,
		cookieSecure: cookieSecure,
		sessionTTL:   sessionTTL,
	}
}

func (h *AuthHandler) setSessionCookie(c *gin.Context, token string) {
	c.SetCookieData(&http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(h.sessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   h.cookieSecure,          // exige HTTPS em producao
		SameSite: http.SameSiteStrictMode, // defesa contra CSRF
	})
}

// clearSessionCookie apaga: mesmo nome e path, MaxAge negativo.
func (h *AuthHandler) clearSessionCookie(c *gin.Context) {
	c.SetCookieData(&http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
}

type credentialsRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// Não existe loginResponse com o token: ele vai no cookie, nunca no corpo.

type meResponse struct {
	UserID string `json:"userId"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req credentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.register.Execute(c.Request.Context(), req.Email, req.Password)
	if errors.Is(err, domain.ErrEmailTaken) {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}
	if err != nil {
		h.log.Error("register failed", zap.Error(err)) // sem e-mail, sem senha
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req credentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.login.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			// 401 genérico — não dizer SE o e-mail existe.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		h.log.Error("login failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	h.setSessionCookie(c, token)
	c.Status(http.StatusNoContent) // 204: sem corpo, o token está no cookie
}

// Logout precisa ser endpoint: JavaScript não consegue apagar cookie HttpOnly.
func (h *AuthHandler) Logout(c *gin.Context) {
	h.clearSessionCookie(c)
	c.Status(http.StatusNoContent)
}

// Me responde se a sessão vale. O front não consegue ler o cookie,
// então precisa perguntar ao servidor depois de recarregar a página.
// Fica no grupo protegido: se o cookie for invalido, o middleware devolve 401.
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := userIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, meResponse{UserID: userID.String()})
}
