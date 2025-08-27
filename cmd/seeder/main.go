package main

import (
	"LojaGin/internal/config"
	"LojaGin/internal/database"
	"LojaGin/internal/seeder"
	"log"
)

func main() {
	log.Println("Iniciando seeder...")

	config.InitEnv()
	db := database.InitDB()

	seeder.Run(db)
}
