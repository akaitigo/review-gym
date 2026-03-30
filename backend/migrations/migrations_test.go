package migrations_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrationFilesExist(t *testing.T) {
	migrationsDir := "."

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("failed to read migrations directory: %v", err)
	}

	sqlFiles := make(map[string]bool)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			sqlFiles[entry.Name()] = true
		}
	}

	if len(sqlFiles) == 0 {
		t.Fatal("no migration SQL files found")
	}

	// Verify migration files come in up/down pairs
	for name := range sqlFiles {
		if strings.HasSuffix(name, ".up.sql") {
			downName := strings.Replace(name, ".up.sql", ".down.sql", 1)
			if !sqlFiles[downName] {
				t.Errorf("migration %s has no corresponding down migration", name)
			}
		}
	}
}

func TestMigrationFilesNotEmpty(t *testing.T) {
	migrationsDir := "."

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("failed to read migrations directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
		if err != nil {
			t.Errorf("failed to read %s: %v", entry.Name(), err)
			continue
		}

		if len(strings.TrimSpace(string(content))) == 0 {
			t.Errorf("migration file %s is empty", entry.Name())
		}
	}
}

func TestInitialSchemaUpContainsTables(t *testing.T) {
	content, err := os.ReadFile("000001_initial_schema.up.sql")
	if err != nil {
		t.Fatalf("failed to read initial schema: %v", err)
	}

	sql := string(content)

	expectedTables := []string{
		"user_profiles",
		"exercises",
		"review_comments",
		"reference_reviews",
		"scores",
	}

	for _, table := range expectedTables {
		if !strings.Contains(sql, "CREATE TABLE "+table) {
			t.Errorf("initial schema missing CREATE TABLE %s", table)
		}
	}
}

func TestInitialSchemaDownDropsTables(t *testing.T) {
	content, err := os.ReadFile("000001_initial_schema.down.sql")
	if err != nil {
		t.Fatalf("failed to read initial schema down: %v", err)
	}

	sql := string(content)

	expectedTables := []string{
		"scores",
		"reference_reviews",
		"review_comments",
		"exercises",
		"user_profiles",
	}

	for _, table := range expectedTables {
		if !strings.Contains(sql, "DROP TABLE IF EXISTS "+table) {
			t.Errorf("down migration missing DROP TABLE IF EXISTS %s", table)
		}
	}
}

func TestUpdateCategoriesMigrationUp(t *testing.T) {
	content, err := os.ReadFile("000002_update_categories.up.sql")
	if err != nil {
		t.Fatalf("failed to read migration 000002 up: %v", err)
	}

	sql := string(content)

	// Verify new category system is defined
	expectedCategories := []string{
		"security",
		"performance",
		"design",
		"readability",
		"error-handling",
	}

	for _, cat := range expectedCategories {
		if !strings.Contains(sql, cat) {
			t.Errorf("migration 000002 up missing category %q", cat)
		}
	}

	// Verify new columns are added
	expectedColumns := []string{
		"category_tags",
		"total_exercises_completed",
		"consecutive_days",
		"last_practice_at",
		"attempt_number",
		"duration_seconds",
	}

	for _, col := range expectedColumns {
		if !strings.Contains(sql, col) {
			t.Errorf("migration 000002 up missing column %q", col)
		}
	}
}

func TestUpdateCategoriesMigrationDown(t *testing.T) {
	content, err := os.ReadFile("000002_update_categories.down.sql")
	if err != nil {
		t.Fatalf("failed to read migration 000002 down: %v", err)
	}

	sql := string(content)

	// Verify rollback restores original categories
	originalCategories := []string{
		"correctness",
		"maintainability",
	}

	for _, cat := range originalCategories {
		if !strings.Contains(sql, cat) {
			t.Errorf("migration 000002 down missing original category %q", cat)
		}
	}

	// Verify columns are dropped
	droppedColumns := []string{
		"category_tags",
		"total_exercises_completed",
		"consecutive_days",
		"last_practice_at",
		"attempt_number",
		"duration_seconds",
	}

	for _, col := range droppedColumns {
		if !strings.Contains(sql, col) {
			t.Errorf("migration 000002 down missing DROP for column %q", col)
		}
	}
}
