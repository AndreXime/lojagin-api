package server

import (
	"LojaGin/internal/user"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbPath = "./db.db"

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
