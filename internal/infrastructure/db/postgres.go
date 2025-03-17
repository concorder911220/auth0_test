package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connStr := os.Getenv("DATABASE_URL")
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
	err = runMigrations()
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func runMigrations() error {
	// Create migrations table if it doesn't exist
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Read migration files from the migrations directory
	migrationDir := "migrations"
	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort files to ensure they run in order (e.g., 001_*, 002_*)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// Apply each migration if not already applied
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Check if migration has already been applied
		var exists bool
		err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM migrations WHERE name = $1)", file.Name()).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", file.Name(), err)
		}
		if exists {
			log.Printf("Skipping already applied migration: %s", file.Name())
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join(migrationDir, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file.Name(), err)
		}

		_, err = DB.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
		}

		_, err = DB.Exec("INSERT INTO migrations (name) VALUES ($1)", file.Name())
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", file.Name(), err)
		}

		log.Printf("Applied migration: %s", file.Name())
	}

	return nil
}
