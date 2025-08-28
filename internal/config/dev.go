//go:build dev

package config

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	ENV_MODE = "dev"
	if err := godotenv.Load(); err != nil {
		log.Println(".env não encontrado, usando defaults/OS")
	}
}
