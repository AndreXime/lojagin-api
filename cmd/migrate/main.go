package main

import (
	"LojaGin/internal/config"
	"LojaGin/internal/database"
	"LojaGin/internal/database/migrations"
	"log"
	"os"
)

func main() {
	config.InitEnv()
	db := database.InitDB()

	m := migrations.GetAllMigrations(db)

	cmd := "up" // Padrão
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "up":
		if err := m.Migrate(); err != nil {
			log.Fatalf("Não foi possível aplicar as migrações: %v", err)
		}
		log.Println("Migrações aplicadas com sucesso.")

	case "down":
		if err := m.RollbackLast(); err != nil {
			log.Fatalf("Não foi possível reverter a última migração: %v", err)
		}
		log.Println("Última migração revertida com sucesso.")

	default:
		log.Fatalf("Comando desconhecido: %s. Use 'up' ou 'down'.", cmd)
	}
}
