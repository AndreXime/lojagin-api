package database

import (
	"LojaGin/internal/config"
	"LojaGin/internal/modules/category"
	"LojaGin/internal/modules/product"
	"LojaGin/internal/modules/user"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.DB_URL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&user.User{}, &category.Category{}, &product.Product{})
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
