package config

import (
	"database/sql"
	"log"

	migrate "github.com/rubenv/sql-migrate"
)

// RunMigrations menjalankan semua migration dari folder /migrations
func RunMigrations(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("Gagal menjalankan migration: %v", err)
	}
	log.Printf("Migration berhasil: %d file dijalankan", n)
}
