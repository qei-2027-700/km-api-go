package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"km-api-go/internal/domain"
	"km-api-go/internal/user/repository/mocks"
)

func TestUserUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUserUsecase(mockRepo)

	tests := []struct {
		name        string
		inputName   string
		inputEmail  string
		inputPass   string
		setupMock   func()
		expectUser  *domain.User
		expectError bool
		errorMsg    string
	}{
		{
			name:       "正常系: 新しいユーザー作成",
			inputName:  "Test User",
			inputEmail: "test@example.com",
			inputPass:  "password123",
			setupMock: func() {
				mockRepo.EXPECT().
					ExistsByEmail("test@example.com").
					Return(false, nil).
					Times(1)
				mockRepo.EXPECT().
					Create(gomock.Any()).
					Do(func(user *domain.User) {
						user.ID = 1 // 作成後のIDを設定
					}).
					Return(nil).
					Times(1)
			},
			expectUser: &domain.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "", // レスポンスではパスワードは除外される
			},
			expectError: false,
		},
		{
			name:       "異常系: メールアドレス重複",
			inputName:  "Test User",
			inputEmail: "duplicate@example.com",
			inputPass:  "password123",
			setupMock: func() {
				mockRepo.EXPECT().
					ExistsByEmail("duplicate@example.com").
					Return(true, nil).
					Times(1)
			},
			expectUser:  nil,
			expectError: true,
			errorMsg:    "user with email duplicate@example.com already exists",
		},
		{
			name:       "異常系: リポジトリエラー（メール存在確認時）",
			inputName:  "Test User",
			inputEmail: "test@example.com",
			inputPass:  "password123",
			setupMock: func() {
				mockRepo.EXPECT().
					ExistsByEmail("test@example.com").
					Return(false, errors.New("database error")).
					Times(1)
			},
			expectUser:  nil,
			expectError: true,
			errorMsg:    "failed to check email existence",
		},
		{
			name:       "異常系: リポジトリエラー（ユーザー作成時）",
			inputName:  "Test User",
			inputEmail: "test@example.com",
			inputPass:  "password123",
			setupMock: func() {
				mockRepo.EXPECT().
					ExistsByEmail("test@example.com").
					Return(false, nil).
					Times(1)
				mockRepo.EXPECT().
					Create(gomock.Any()).
					Return(errors.New("database error")).
					Times(1)
			},
			expectUser:  nil,
			expectError: true,
			errorMsg:    "failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			user, err := usecase.Create(context.Background(), tt.inputName, tt.inputEmail, tt.inputPass)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectUser.ID, user.ID)
				assert.Equal(t, tt.expectUser.Name, user.Name)
				assert.Equal(t, tt.expectUser.Email, user.Email)
				assert.Equal(t, "", user.Password) // パスワードは除外されている
			}
		})
	}
}

func TestUserUsecase_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUserUsecase(mockRepo)

	tests := []struct {
		name        string
		setupMock   func()
		expectUsers []domain.User
		expectError bool
	}{
		{
			name: "正常系: ユーザー一覧取得",
			setupMock: func() {
				users := []domain.User{
					{ID: 1, Name: "User1", Email: "user1@example.com", Password: "hashed1"},
					{ID: 2, Name: "User2", Email: "user2@example.com", Password: "hashed2"},
				}
				mockRepo.EXPECT().GetAll().Return(users, nil).Times(1)
			},
			expectUsers: []domain.User{
				{ID: 1, Name: "User1", Email: "user1@example.com", Password: ""},
				{ID: 2, Name: "User2", Email: "user2@example.com", Password: ""},
			},
			expectError: false,
		},
		{
			name: "正常系: 空のユーザー一覧",
			setupMock: func() {
				mockRepo.EXPECT().GetAll().Return([]domain.User{}, nil).Times(1)
			},
			expectUsers: []domain.User{},
			expectError: false,
		},
		{
			name: "異常系: リポジトリエラー",
			setupMock: func() {
				mockRepo.EXPECT().GetAll().Return(nil, errors.New("database error")).Times(1)
			},
			expectUsers: nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			users, err := usecase.GetAllUsers(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, users)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectUsers, users)
				// パスワードが除外されていることを確認
				for _, user := range users {
					assert.Equal(t, "", user.Password)
				}
			}
		})
	}
}

func TestUserUsecase_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUserUsecase(mockRepo)

	tests := []struct {
		name        string
		inputID     uint
		setupMock   func()
		expectUser  *domain.User
		expectError bool
	}{
		{
			name:    "正常系: ユーザー取得",
			inputID: 1,
			setupMock: func() {
				user := &domain.User{
					ID:       1,
					Name:     "Test User",
					Email:    "test@example.com",
					Password: "hashedpassword",
				}
				mockRepo.EXPECT().GetByID(uint(1)).Return(user, nil).Times(1)
			},
			expectUser: &domain.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "", // パスワードは除外される
			},
			expectError: false,
		},
		{
			name:    "異常系: ユーザーが見つからない",
			inputID: 999,
			setupMock: func() {
				mockRepo.EXPECT().GetByID(uint(999)).Return(nil, errors.New("user not found")).Times(1)
			},
			expectUser:  nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			user, err := usecase.GetUserByID(context.Background(), tt.inputID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectUser.ID, user.ID)
				assert.Equal(t, tt.expectUser.Name, user.Name)
				assert.Equal(t, tt.expectUser.Email, user.Email)
				assert.Equal(t, "", user.Password)
			}
		})
	}
}

func TestUserUsecase_AuthenticateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUserUsecase(mockRepo)

	// テスト用のパスワードハッシュを生成
	testUser := &domain.User{
		Password: "password123",
	}
	err := testUser.HashPassword()
	assert.NoError(t, err)
	hashedPassword := testUser.Password

	tests := []struct {
		name        string
		inputEmail  string
		inputPass   string
		setupMock   func()
		expectUser  *domain.User
		expectError bool
	}{
		{
			name:       "正常系: 認証成功",
			inputEmail: "test@example.com",
			inputPass:  "password123",
			setupMock: func() {
				user := &domain.User{
					ID:       1,
					Name:     "Test User",
					Email:    "test@example.com",
					Password: hashedPassword, // 実際のハッシュ化されたパスワードを使用
				}
				mockRepo.EXPECT().GetByEmail("test@example.com").Return(user, nil).Times(1)
			},
			expectUser: &domain.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "",
			},
			expectError: false,
		},
		{
			name:       "異常系: パスワード不一致",
			inputEmail: "test@example.com",
			inputPass:  "wrongpassword",
			setupMock: func() {
				user := &domain.User{
					ID:       1,
					Name:     "Test User",
					Email:    "test@example.com",
					Password: hashedPassword,
				}
				mockRepo.EXPECT().GetByEmail("test@example.com").Return(user, nil).Times(1)
			},
			expectUser:  nil,
			expectError: true,
		},
		{
			name:       "異常系: ユーザーが見つからない",
			inputEmail: "notfound@example.com",
			inputPass:  "password123",
			setupMock: func() {
				mockRepo.EXPECT().GetByEmail("notfound@example.com").Return(nil, errors.New("user not found")).Times(1)
			},
			expectUser:  nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			user, err := usecase.AuthenticateUser(context.Background(), tt.inputEmail, tt.inputPass)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "authentication failed")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectUser.ID, user.ID)
				assert.Equal(t, tt.expectUser.Name, user.Name)
				assert.Equal(t, tt.expectUser.Email, user.Email)
				assert.Equal(t, "", user.Password)
			}
		})
	}
}