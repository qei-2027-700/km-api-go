package user

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"km-api-go/internal/helper"
)

// UserHandler はユーザー関連のHTTPリクエストを処理します
type UserHandler struct {
	usecase UserUsecase
}

// NewUserHandler は新しいUserHandlerを初期化します
func NewUserHandler(usecase UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

// CreateUser は新しいユーザーを作成します
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, helper.ErrorCodeValidation, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return helper.ValidationErrorResponse(c, err.Error())
	}

	user, err := h.usecase.Create(c.Request().Context(), req.Name, req.Email, req.Password)
	if err != nil {
		// TODO: エラーの種類に応じて、AlreadyExistsResponseなども返す
		return helper.InternalErrorResponse(c, err.Error())
	}

	res := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return helper.SuccessResponse(c, http.StatusCreated, res, "User created successfully")
}

// TODO: GetUser, UpdateUser, DeleteUser などのハンドラーメソッドを実装
