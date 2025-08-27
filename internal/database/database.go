package database

import (
	"LojaGin/internal/config"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.DB_URL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
