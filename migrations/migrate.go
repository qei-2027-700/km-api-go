package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"km-api-go/internal/infra"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// .envファイル読み込み
	if err := godotenv.Overload(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// データベース接続
	db, err := infra.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := infra.CloseDatabase(db); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	// マイグレーション実行
	if err := runMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("✅ All migrations completed successfully!")
}

func runMigrations(db interface{}) error {
	migrationsDir := "migrations"
	
	// マイグレーションファイルを取得
	var migrationFiles []string
	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			migrationFiles = append(migrationFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// ファイル名でソート
	sort.Strings(migrationFiles)

	log.Printf("Found %d migration files", len(migrationFiles))

	// 各マイグレーションファイルを実行
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		return fmt.Errorf("db is not a *gorm.DB")
	}

	for _, file := range migrationFiles {
		log.Printf("Running migration: %s", file)

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// SQLを実行
		if err := gormDB.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute migration file %s: %w", file, err)
		}
	}

	return nil
}