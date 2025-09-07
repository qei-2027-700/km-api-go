package helper

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// APIResponse 統一APIレスポンス形式
// @Description 統一されたAPIレスポンス形式
type APIResponse struct {
	Success bool        `json:"success" example:"true"`                    // 成功フラグ
	Data    interface{} `json:"data,omitempty"`                           // レスポンスデータ
	Message string      `json:"message,omitempty" example:"操作が成功しました"`      // メッセージ
	Error   *APIError   `json:"error,omitempty"`                          // エラー情報
}

// APIError APIエラー情報
// @Description APIエラーの詳細情報
type APIError struct {
	Code    ErrorCode `json:"code" example:"VALIDATION_ERROR"`       // エラーコード
	Message string    `json:"message" example:"バリデーションエラーが発生しました"` // エラーメッセージ
	Details string    `json:"details,omitempty"`                     // エラー詳細
}

// PaginatedResponse ページネーション付きレスポンス
// @Description ページネーション付きのレスポンス形式
type PaginatedResponse struct {
	Success    bool                `json:"success" example:"true"`    // 成功フラグ
	Data       interface{}         `json:"data"`                      // データ配列
	Pagination *PaginationResponse `json:"pagination"`                // ページネーション情報
	Message    string              `json:"message,omitempty"`         // メッセージ
	Error      *APIError           `json:"error,omitempty"`           // エラー情報
}

// SuccessResponse 成功レスポンスを作成
func SuccessResponse(c echo.Context, statusCode int, data interface{}, message string) error {
	return c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorResponse エラーレスポンスを作成
func ErrorResponse(c echo.Context, statusCode int, code ErrorCode, message string, details string) error {
	return c.JSON(statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// PaginatedSuccessResponse ページネーション付き成功レスポンスを作成
func PaginatedSuccessResponse(c echo.Context, data interface{}, pagination *PaginationResponse, message string) error {
	return c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
		Message:    message,
	})
}

// CreatedResponse 作成成功レスポンス
func CreatedResponse(c echo.Context, data interface{}, message string) error {
	if message == "" {
		message = "リソースが正常に作成されました"
	}
	return SuccessResponse(c, http.StatusCreated, data, message)
}

// UpdatedResponse 更新成功レスポンス
func UpdatedResponse(c echo.Context, data interface{}, message string) error {
	if message == "" {
		message = "リソースが正常に更新されました"
	}
	return SuccessResponse(c, http.StatusOK, data, message)
}

// DeletedResponse 削除成功レスポンス
func DeletedResponse(c echo.Context, message string) error {
	if message == "" {
		message = "リソースが正常に削除されました"
	}
	return SuccessResponse(c, http.StatusOK, nil, message)
}

// NotFoundResponse 404エラーレスポンス
func NotFoundResponse(c echo.Context, resource string) error {
	message := "リソースが見つかりません"
	if resource != "" {
		message = resource + "が見つかりません"
	}
	return ErrorResponse(c, http.StatusNotFound, ErrorCodeNotFound, message, "")
}

// ValidationErrorResponse バリデーションエラーレスポンス
func ValidationErrorResponse(c echo.Context, details string) error {
	return ErrorResponse(c, http.StatusBadRequest, ErrorCodeValidation, "入力データが正しくありません", details)
}

// AlreadyExistsResponse 既存リソースエラーレスポンス
func AlreadyExistsResponse(c echo.Context, resource string) error {
	message := "リソースは既に存在します"
	if resource != "" {
		message = resource + "は既に存在します"
	}
	return ErrorResponse(c, http.StatusConflict, ErrorCodeAlreadyExists, message, "")
}

// UnauthorizedResponse 認証エラーレスポンス
func UnauthorizedResponse(c echo.Context) error {
	return ErrorResponse(c, http.StatusUnauthorized, ErrorCodeUnauthorized, "認証が必要です", "")
}

// ForbiddenResponse 認可エラーレスポンス
func ForbiddenResponse(c echo.Context) error {
	return ErrorResponse(c, http.StatusForbidden, ErrorCodeForbidden, "このリソースにアクセスする権限がありません", "")
}

// InternalErrorResponse 内部エラーレスポンス
func InternalErrorResponse(c echo.Context, details string) error {
	return ErrorResponse(c, http.StatusInternalServerError, ErrorCodeInternalError, "内部サーバーエラーが発生しました", details)
}

// DatabaseErrorResponse データベースエラーレスポンス
func DatabaseErrorResponse(c echo.Context, details string) error {
	return ErrorResponse(c, http.StatusInternalServerError, ErrorCodeDatabaseError, "データベースエラーが発生しました", details)
}