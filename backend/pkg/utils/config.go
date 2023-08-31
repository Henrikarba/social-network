package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error getting env key: %v", err)
	}
	return os.Getenv(key)
}
