package tests

import (
	"LojaGin/internal/config"
	"LojaGin/internal/database"
	"LojaGin/internal/routes"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine
var db *gorm.DB

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// 2. Define um nome para o banco de dados de teste.
	config.DB_URL = "test_e2e.db"
	config.JWT_SECRET = "segredo_para_testes_123456"

	// 3. Garante que qualquer arquivo de DB de um teste anterior seja removido ANTES de começar.
	// Isso garante um estado 100% limpo.
	os.Remove(config.DB_URL)

	// 4. Agora, inicializa uma nova conexão de banco de dados limpa.
	db = database.InitDB()

	// 5. Cria o router e injeta a conexão do DB.
	router = gin.Default()
	routes.SetupAPI(router, db)

	// 6. Roda os testes.
	exitCode := m.Run()

	// 7. Limpeza final após todos os testes.
	os.Remove(config.DB_URL)
	os.Exit(exitCode)
}

// clearDB limpa todas as tabelas para garantir que os testes sejam independentes
func clearDB() {
	// Deleta os registros e reseta a sequência de auto incremento do SQLite
	db.Exec("DELETE FROM order_items")
	db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'order_items'")
	db.Exec("DELETE FROM orders")
	db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'orders'")
	db.Exec("DELETE FROM cart_items")
	db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'cart_items'")
	db.Exec("DELETE FROM carts")
	db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'carts'")
	db.Exec("DELETE FROM products")
	db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'products'")
	db.Exec("DELETE FROM categories")
	db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'categories'")
	db.Exec("DELETE FROM users")
	db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'users'")
}
