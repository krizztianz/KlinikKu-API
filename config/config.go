package config

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DB         *sql.DB
	ServerPort string
	JWTSecret  string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // jika gagal, gunakan env dari OS

	dbURL := os.Getenv("DATABASE_URL")
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	cfg := &Config{
		DB:         conn,
		ServerPort: os.Getenv("PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	return cfg, nil
}
