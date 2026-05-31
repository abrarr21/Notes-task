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

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file %v", err)
	}

	mongo_uri := os.Getenv("MONGODB_URI")
	if mongo_uri == "" {
		log.Fatal("No MongoDB_URI defined in the .env file")
	}

	dbname := GetEnv("DB_NAME", "")
	if dbname == "" {
		log.Fatal("No DB_NAME is provided in the .env file")
	}

	return &Config{
		ServerConfig{
			Port: GetEnv("PORT", "8080"),
			Env:  GetEnv("ENV", "development"),
		},

		DatabaseConfig{
			MongoDB_URI: mongo_uri,
			DBName:      GetEnv("DB_Name", "my-app-dev"),
		},
	}
}

func GetEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}
