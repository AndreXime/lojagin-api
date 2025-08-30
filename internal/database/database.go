package database

import (
	"LojaGin/internal/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	var db *gorm.DB
	var err error
	const maxRetries = 5
	const baseDelay = 1 * time.Second

	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(config.DB_URL), &gorm.Config{})
		if err == nil {
			log.Println("Conexão com o banco de dados estabelecida com sucesso!")
			return db
		}

		// Atraso exponencial: o tempo de espera aumenta a cada tentativa.
		// Isso evita sobrecarregar o banco de dados durante uma falha prolongada ou inicialização.
		delay := baseDelay * time.Duration(1<<uint(i-1))

		log.Printf("Tentativa %d de %d falhou. Tentando em %v segundos...", i, maxRetries, delay)

		time.Sleep(delay)
	}

	log.Fatalln("Não foi possível conectar ao banco de dados após 10 tentativas. Abortando.")
	return nil // Nunca será alcançado devido ao log.Fatalln
}
