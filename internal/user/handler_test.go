package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"km-api-go/internal/domain"
	"km-api-go/internal/helper"
	"km-api-go/internal/user/mocks"
)

func TestUserHandler_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	handler := NewUserHandler(mockUsecase)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, responseBody string)
	}{
		{
			name: "正常系: ユーザー作成成功",
			requestBody: CreateUserRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				user := &domain.User{
					ID:    1,
					Name:  "Test User",
					Email: "test@example.com",
				}
				mockUsecase.EXPECT().
					Create(gomock.Any(), "Test User", "test@example.com", "password123").
					Return(user, nil).
					Times(1)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, responseBody string) {
				var response helper.APIResponse
				err := json.Unmarshal([]byte(responseBody), &response)
				assert.NoError(t, err)
				assert.True(t, response.Success)
				assert.Equal(t, "User created successfully", response.Message)
				assert.NotNil(t, response.Data)

				// データの詳細をチェック
				userData, ok := response.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, float64(1), userData["id"])
				assert.Equal(t, "Test User", userData["name"])
				assert.Equal(t, "test@example.com", userData["email"])
			},
		},
		{
			name: "異常系: 不正なリクエストボディ",
			requestBody: map[string]interface{}{
				"name":     123, // 不正な型
				"email":    "test@example.com",
				"password": "password123",
			},
			setupMock:      func() {}, // モック呼び出しなし
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, responseBody string) {
				var response helper.APIResponse
				err := json.Unmarshal([]byte(responseBody), &response)
				assert.NoError(t, err)
				assert.False(t, response.Success)
				assert.NotNil(t, response.Error)
			},
		},
		{
			name: "異常系: バリデーションエラー（空のフィールド）",
			requestBody: CreateUserRequest{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock:      func() {}, // バリデーションで失敗するため、ユースケースは呼ばれない
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, responseBody string) {
				var response helper.APIResponse
				err := json.Unmarshal([]byte(responseBody), &response)
				assert.NoError(t, err)
				assert.False(t, response.Success)
			},
		},
		{
			name: "異常系: ユースケースでエラー発生",
			requestBody: CreateUserRequest{
				Name:     "Test User",
				Email:    "duplicate@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockUsecase.EXPECT().
					Create(gomock.Any(), "Test User", "duplicate@example.com", "password123").
					Return(nil, errors.New("user with email duplicate@example.com already exists")).
					Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, responseBody string) {
				var response helper.APIResponse
				err := json.Unmarshal([]byte(responseBody), &response)
				assert.NoError(t, err)
				assert.False(t, response.Success)
				assert.NotNil(t, response.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			e.Validator = helper.NewValidator()

			// リクエストボディをJSON化
			reqBodyBytes, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(reqBodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tt.setupMock()

			// テスト実行
			err = handler.CreateUser(c)

			// アサーション
			if assert.NoError(t, err) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				tt.checkResponse(t, rec.Body.String())
			}
		})
	}
}