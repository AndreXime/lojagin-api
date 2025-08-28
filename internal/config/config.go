package config

import "os"

var (
	JWT_SECRET = "DNOASBDOABD@#!&*(#@!&#(*))"
	DB_URL     = "db.db"
	PORT       = "8080"

	// ENV_MODE pode ser mudado com tag na build dev
	ENV_MODE = "production"
)

func InitEnv() {
	if v := os.Getenv("JWT_SECRET"); v != "" {
		JWT_SECRET = v
	}
	if v := os.Getenv("DB_URL"); v != "" {
		DB_URL = v
	}
	if v := os.Getenv("PORT"); v != "" {
		PORT = v
	}
}
