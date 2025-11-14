package utils

import (
	"os"

	"github.com/joho/godotenv"
)

// Config func to get env value
func Env(key string) string {
	// load .env file
	godotenv.Load(".env")
	return os.Getenv(key)
}
