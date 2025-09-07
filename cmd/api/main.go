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
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"km-api-go/internal/helper"
	"km-api-go/internal/infra"
	"km-api-go/internal/user"
	userRepo "km-api-go/internal/user/repository"
)

func main() {
	// .envファイル読み込み
	if err := godotenv.Overload(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Echoインスタンス生成
	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// カスタムバリデータ設定
	e.Validator = helper.NewValidator()

	// データベース接続
	db, err := infra.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 依存関係の注入 (Dependency Injection)
	userRepository := userRepo.NewUserRepository(db)
	userUsecase := user.NewUserUsecase(userRepository)
	userHandler := user.NewUserHandler(userUsecase)

	// ルーティング
	apiV1 := e.Group("/api/v1")
	apiV1.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	usersGroup := apiV1.Group("/users")
	usersGroup.POST("", userHandler.CreateUser)

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
