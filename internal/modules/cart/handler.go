package cart

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// getUserIDFromContext é uma função helper para extrair o ID do usuário do contexto do Gin
func getUserIDFromContext(c *gin.Context) (uint, error) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("ID do usuário não encontrado no contexto")
	}
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		return 0, errors.New("ID do usuário inválido no contexto")
	}
	return uint(userID), nil
}

func (h *Handler) GetCart(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Acesso não autorizado"})
		return
	}
	cart, err := h.service.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (h *Handler) AddToCart(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Acesso não autorizado"})
		return
	}
	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	cart, err := h.service.AddToCart(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (h *Handler) RemoveFromCart(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Acesso não autorizado"})
		return
	}
	var req RemoveFromCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	cart, err := h.service.RemoveFromCart(userID, req)
	if err != nil {
		if err == ErrProductNotInCart || err == ErrInsufficientStock {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (h *Handler) ClearCart(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Acesso não autorizado"})
		return
	}
	if err := h.service.ClearCart(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) Checkout(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Acesso não autorizado"})
		return
	}
	order, err := h.service.Checkout(userID)
	if err != nil {
		if err == ErrCartNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Não é possível fazer checkout de um carrinho vazio"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}
