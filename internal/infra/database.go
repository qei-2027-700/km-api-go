package infra

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig データベース接続設定
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadDatabaseConfig 環境変数からデータベース設定を読み込み
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "km_api"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// getEnv 環境変数取得のヘルパー関数
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewDatabase データベース接続を初期化
func NewDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	// PostgreSQL DSN構築
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Tokyo",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)

	// ログレベル設定
	logLevel := logger.Info
	if getEnv("GO_ENV", "development") == "production" {
		logLevel = logger.Error
	}

	// GORM設定
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// データベース接続
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 接続プール設定
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 接続プール設定
	sqlDB.SetMaxIdleConns(10)                // アイドル接続数
	sqlDB.SetMaxOpenConns(100)               // 最大接続数
	sqlDB.SetConnMaxLifetime(time.Hour * 1)  // 接続最大生存時間

	log.Println("Database connected successfully")
	return db, nil
}

// InitDatabase データベース初期化（環境変数から設定読み込み）
func InitDatabase() (*gorm.DB, error) {
	config := LoadDatabaseConfig()
	return NewDatabase(config)
}

// CloseDatabase データベース接続を閉じる
func CloseDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance for closing: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("Database connection closed")
	return nil
}

// PingDatabase データベース接続確認
func PingDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}