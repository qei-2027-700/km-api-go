package user

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"km-api-go/internal/helper"
)

type UserHandler struct {
	usecase UserUsecase
}

func NewUserHandler(usecase UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

// CreateUser godoc
// @Summary ユーザー作成
// @Description 新しいユーザーを作成します
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "ユーザー情報"
// @Success 201 {object} helper.APIResponse{data=UserResponse}
// @Failure 400 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /users [post]
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
