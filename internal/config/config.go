package config

import "os"

var JWT_SECRET string
var DB_URL string
var ENV_MODE string

func InitEnv() {
	value := os.Getenv("JWT_SECRET")
	if value == "" {
		JWT_SECRET = "DNOASBDOABD@#!&*(#@!&#(*))"
	} else {
		JWT_SECRET = value
	}

	value = os.Getenv("DB_URL")
	if value == "" {
		DB_URL = "db.db"
	} else {
		DB_URL = value
	}

	value = os.Getenv("ENV_MODE")
	if value == "" {
		ENV_MODE = "dev"
	} else {
		ENV_MODE = value
	}
}
