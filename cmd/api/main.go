// Package main KM API Go
//
// @title KM API
// @version 1.0
// @description KM API for Go backend application
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /api/v1
//
// @schemes http https
//
// @tag.name users
// @tag.description ユーザー関連のAPI
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"km-api-go/internal/infra"
	"km-api-go/server"
	
	// Swagger docs
	_ "km-api-go/swagger/src"
)

func main() {
	// .envファイル読み込み
	if err := godotenv.Overload(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// データベース接続
	db, err := infra.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// ルーターのセットアップ
	e := server.SetupRouter(db)

	// サーバーの起動とグレースフルシャットダウン
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	log.Println("Server gracefully stopped")

	if err := infra.CloseDatabase(db); err != nil {
		log.Printf("Failed to close database: %v", err)
	}

	log.Println("Database connection closed")
}
