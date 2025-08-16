package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type AppClaims struct {
	jwt.RegisteredClaims
}

// Middleware para proteger rotas que requerem autenticação.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Cabeçalho de autorização não encontrado"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato do cabeçalho de autorização inválido"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &AppClaims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("método de assinatura inesperado")
			}
			return jwtSecret, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			}
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*AppClaims); ok && token.Valid {
			c.Set("userID", claims.Subject)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
		}
	}
}
