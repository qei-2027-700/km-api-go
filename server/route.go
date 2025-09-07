package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"

	"km-api-go/internal/helper"
	"km-api-go/internal/user"
	userRepo "km-api-go/internal/user/repository"
)

func SetupRouter(db *gorm.DB) *echo.Echo {
	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// カスタムバリデータ設定
	e.Validator = helper.NewValidator()

	// 依存関係の注入 (Dependency Injection)
	userRepository := userRepo.NewUserRepository(db)
	userUsecase := user.NewUserUsecase(userRepository)
	userHandler := user.NewUserHandler(userUsecase)

	// ルーティング
	apiV1 := e.Group("/api/v1")

	// ヘルスチェック
	apiV1.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// ユーザー関連
	usersGroup := apiV1.Group("/users")
	usersGroup.POST("", userHandler.CreateUser)

	return e
}
