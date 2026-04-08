package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	direction := flag.String("direction", "up", "migration direction: up or down")
	migrationsPath := flag.String("path", "file://migrations", "path to migrations directory")
	flag.Parse()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://review_gym:review_gym_dev@localhost:5432/review_gym?sslmode=disable"
	}

	// golang-migrate's pgx/v5 driver requires the "pgx5://" URL scheme.
	// Convert common PostgreSQL URL schemes for compatibility.
	migrateURL := databaseURL
	if strings.HasPrefix(migrateURL, "postgresql://") {
		migrateURL = "pgx5://" + strings.TrimPrefix(migrateURL, "postgresql://")
	} else if strings.HasPrefix(migrateURL, "postgres://") {
		migrateURL = "pgx5://" + strings.TrimPrefix(migrateURL, "postgres://")
	}

	m, err := migrate.New(*migrationsPath, migrateURL)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("source close error: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("database close error: %v", dbErr)
		}
	}()

	switch *direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up failed: %v", err)
		}
		fmt.Println("migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration down failed: %v", err)
		}
		fmt.Println("migrations rolled back successfully")
	default:
		log.Fatalf("unknown direction: %s (use 'up' or 'down')", *direction)
		os.Exit(1)
	}
}
