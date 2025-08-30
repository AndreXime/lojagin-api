package config

import (
	"log"
	"os"
)

var (
	PORT       = "8080"
	JWT_SECRET string
	DB_URL     string

	// ENV_MODE pode ser mudado com tag na build dev
	ENV_MODE = "production"
)

func InitEnv() {
	v := os.Getenv("JWT_SECRET")
	if v != "" {
		JWT_SECRET = v
	} else {
		log.Fatalln("Não foi encontrado a variavel JWT_SECRET")
	}

	v = os.Getenv("DB_URL")
	if v != "" {
		DB_URL = v
	} else {
		log.Fatalln("Não foi encontrado a variavel DB_URL")
	}

	if v := os.Getenv("PORT"); v != "" {
		PORT = v
	}
}
