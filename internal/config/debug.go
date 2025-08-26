package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

type routeWithParts struct {
	Method  string
	Path    string
	Handler string
	Parts   []string
}

func PrintRoutes(router *gin.Engine) {
	color.Cyan("Rotas registradas:")

	routes := router.Routes()

	// pré-processa paths em segmentos uma vez só
	processed := make([]routeWithParts, len(routes))
	for i, r := range routes {
		processed[i] = routeWithParts{
			Method:  r.Method,
			Path:    r.Path,
			Handler: r.Handler,
			Parts:   strings.Split(strings.Trim(r.Path, "/"), "/"),
		}
	}

	// ordena comparando só os slices prontos
	sort.SliceStable(processed, func(i, j int) bool {
		pi, pj := processed[i].Parts, processed[j].Parts
		for k := 0; k < len(pi) && k < len(pj); k++ {
			if pi[k] != pj[k] {
				return pi[k] < pj[k]
			}
		}
		return len(pi) < len(pj)
	})

	for _, ri := range processed {
		var methodColor *color.Color
		switch ri.Method {
		case "GET":
			methodColor = color.New(color.FgGreen, color.Bold)
		case "POST":
			methodColor = color.New(color.FgBlue, color.Bold)
		case "PUT":
			methodColor = color.New(color.FgYellow, color.Bold)
		case "DELETE":
			methodColor = color.New(color.FgRed, color.Bold)
		default:
			methodColor = color.New(color.FgWhite, color.Bold)
		}

		fmt.Printf("%s %s --> %s\n",
			methodColor.Sprint(fmt.Sprintf("%-7s", ri.Method)),
			color.MagentaString(fmt.Sprintf("%-25s", ri.Path)),
			color.WhiteString(ri.Handler),
		)
	}
}
