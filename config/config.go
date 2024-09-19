package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	DBDriver  string
	DBName    string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	SSLMode   string
	JWTSecret string
	MONGOHost string
	MONGOPort string
	MONGODB   string
	TestDBURL string
)

func init() {
	loadEnv()

	DBDriver = GetEnv("DB_DRIVER")
	DBName = GetEnv("DB_NAME")
	DBHost = GetEnv("DB_HOST")
	DBPort = GetEnv("DB_PORT")
	DBUser = GetEnv("DB_USER")
	DBPass = GetEnv("DB_PASS")
	SSLMode = GetEnv("SSL_MODE")
	JWTSecret = GetEnv("JWT_SECRET")
	MONGOHost = GetEnv("MONGO_HOST")
	MONGOPort = GetEnv("MONGO_PORT")
	MONGODB = GetEnv("MONGO_DB")
	TestDBURL = GetEnv("TEST_DB_URL")
}

func loadEnv() {
	possiblePaths := []string{
		".env",
		filepath.Join("..", ".env"),
		filepath.Join("../..", ".env"),
	}

	var err error
	for _, path := range possiblePaths {
		err = godotenv.Load(path)
		if err == nil {
			log.Printf("Loaded .env file from %s", path)
			return
		}
	}

	log.Printf("Error loading .env file: %v", err)
}

func GetEnv(key string, value ...string) string {
	if envValue := os.Getenv(key); envValue != "" {
		return envValue
	}
	if len(value) > 0 {
		return value[0]
	}
	return ""
}
