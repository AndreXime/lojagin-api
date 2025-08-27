package migrations

import (
	"LojaGin/internal/modules/cart"
	"LojaGin/internal/modules/category"
	"LojaGin/internal/modules/product"
	"LojaGin/internal/modules/user"
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// GetAllMigrations retorna uma lista de todas as migrações do projeto.
func GetAllMigrations(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20250826_initial_schema",
			Migrate: func(tx *gorm.DB) error {
				log.Println("Executando migração: 20250826_initial_schema")
				return tx.AutoMigrate(
					&user.User{},
					&category.Category{},
					&product.Product{},
					&cart.Cart{},
					&cart.CartItem{},
					&cart.Order{},
					&cart.OrderItem{},
				)
			},
			Rollback: func(tx *gorm.DB) error {
				log.Println("Revertendo migração: 20250826_initial_schema")
				return tx.Migrator().DropTable(
					"order_items",
					"orders",
					"cart_items",
					"carts",
					"products",
					"categories",
					"users",
				)
			},
		},
	})
}
