package tests

import (
	"LojaGin/internal/config"
	"LojaGin/internal/database"
	"LojaGin/internal/routes"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine
var db *gorm.DB

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	config.InitEnv()
	db = database.InitDB()

	if err := clearDB(); err != nil {
		log.Fatalln("Erro ao limpar DB:", err)
	}

	router = gin.Default()
	routes.SetupAPI(router, db)

	exitCode := m.Run()

	if err := clearDB(); err != nil {
		log.Println("Erro ao limpar DB após testes:", err)
	}

	os.Exit(exitCode)
}

// clearDB limpa todas as tabelas para garantir que os testes sejam independentes
func clearDB() error {
	tables := []string{
		"order_items",
		"orders",
		"cart_items",
		"carts",
		"products",
		"categories",
		"users",
	}

	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error; err != nil {
			return err
		}
	}
	return nil
}
