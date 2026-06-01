package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	MongoDB_URI string
	DBName      string
}

type JWTConfig struct {
	Secret string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found")
	}

	mongodb_uri := os.Getenv("MONGODB_URI")
	if mongodb_uri == "" {
		log.Fatal("MONGODB_URI is not provided in the .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not provided in the .env file")
	}

	return &Config{
		ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},

		DatabaseConfig{
			MongoDB_URI: mongodb_uri,
			DBName:      getEnv("DB_NAME", "test-dev"),
		},

		JWTConfig{
			Secret: jwtSecret,
		},
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}
