package config

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadConfig() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using default environment variables")
	}
}
