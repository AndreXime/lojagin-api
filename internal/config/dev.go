//go:build dev

package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func init() {
	ENV_MODE = "dev"

	// Sempre vai encontrar .env na raiz do projeto independente quem chamou
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "../..")
	envPath := filepath.Join(projectRoot, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Println(".env não encontrado, usando defaults/OS")
		return
	}

	log.Println("Arquivo .env encontrado")
}
