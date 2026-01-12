package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() Config {
	godotenv.Load()

	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"), 
	}
}