package user

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrEmailExists        = errors.New("o e-mail informado já está em uso")
	ErrUserNotFound       = errors.New("usuário não encontrado")
	ErrInvalidPassword    = errors.New("a senha não pode estar em branco ou conter apenas espaços")
	ErrShortPassword      = errors.New("a senha deve ter no mínimo 8 caracteres")
	ErrLongPassword       = errors.New("a senha não pode ter mais de 72 caracteres")
	ErrInvalidRequestData = errors.New("dados da requisição inválidos ou mal formatados")

	ErrDatabase = errors.New("ocorreu um erro ao processar sua solicitação")
)

func (h *Handler) handleError(c *gin.Context, err error) {
	switch err {
	case ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case ErrEmailExists:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case ErrInvalidRequestData, ErrInvalidPassword, ErrShortPassword, ErrLongPassword:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case ErrDatabase:
		// Loga o erro real para a equipe de desenvolvimento
		log.Printf("Internal database error: %v", err)
		// Retorna uma mensagem genérica para o cliente
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		// Erros inesperados
		log.Printf("Unexpected internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro inesperado"})
	}
}
