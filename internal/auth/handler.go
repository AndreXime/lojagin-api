package auth

import (
	"LojaGin/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Register(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar usuário"})
		return
	}

	c.SetCookie("token", token, 3600*24, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, gin.H{"mensagem": "Usuário registrado com sucesso"})
}

func (h *Handler) Login(c *gin.Context) {
	var req user.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"mensagem": "Login feito com sucesso"})
}
