//go:build !dev

package docs

import (
	"github.com/gin-gonic/gin"
)

// Não existe swagger em produção e para n ter erro de compilação precisa da função setupdocs em package docs
func SetupDocs(router *gin.Engine) {}
