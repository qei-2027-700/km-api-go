package user

import (
	"context"
	"fmt"

	"km-api-go/internal/domain"
	"km-api-go/internal/helper"
	"km-api-go/internal/user/repository"
)

// UserUsecase defines the interface for user business logic.
// It's good practice to depend on interfaces, not concrete implementations.
type UserUsecase interface {
	Create(ctx context.Context, name, email, password string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	GetUserByID(ctx context.Context, id uint) (*domain.User, error)
	UpdateUser(ctx context.Context, id uint, name, email string) (*domain.User, error)
	DeleteUser(ctx context.Context, id uint) error
	GetUsersPaginated(ctx context.Context, page, limit int) ([]domain.User, *helper.PaginationResponse, error)
	AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error)
}

// userUsecase implements the UserUsecase interface.
type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase is the constructor for userUsecase.
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// Create creates a new user.
func (uc *userUsecase) Create(ctx context.Context, name, email, password string) (*domain.User, error) {
	// メール重複チェック
	exists, err := uc.userRepo.ExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// ユーザーオブジェクト作成
	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: password, // BeforeCreateフックでハッシュ化される
	}

	// リポジトリで保存
	if err := uc.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// パスワードを除外したレスポンス用データを返す
	user.Password = ""
	return user, nil
}

// GetAllUsers retrieves all users.
func (uc *userUsecase) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	users, err := uc.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	// パスワードを除外
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

// GetUserByID retrieves a user by their ID.
func (uc *userUsecase) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id %d: %w", id, err)
	}

	user.Password = ""
	return user, nil
}

// UpdateUser updates a user's information.
func (uc *userUsecase) UpdateUser(ctx context.Context, id uint, name, email string) (*domain.User, error) {
	existingUser, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user for update: %w", err)
	}

	if existingUser.Email != email {
		exists, err := uc.userRepo.ExistsByEmail(email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("user with email %s already exists", email)
		}
	}

	existingUser.Name = name
	existingUser.Email = email

	if err := uc.userRepo.Update(existingUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	existingUser.Password = ""
	return existingUser, nil
}

// DeleteUser deletes a user by their ID.
func (uc *userUsecase) DeleteUser(ctx context.Context, id uint) error {
	return uc.userRepo.Delete(id)
}

// GetUsersPaginated retrieves users with pagination.
func (uc *userUsecase) GetUsersPaginated(ctx context.Context, page, limit int) ([]domain.User, *helper.PaginationResponse, error) {
	paginationReq := &helper.PaginationRequest{Page: page, Limit: limit}
	offset := paginationReq.GetOffset()
	normalizedLimit := paginationReq.GetLimit()

	total, err := uc.userRepo.Count()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count users: %w", err)
	}

	users, err := uc.userRepo.GetPaginated(offset, normalizedLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get paginated users: %w", err)
	}

	for i := range users {
		users[i].Password = ""
	}

	pagination := helper.NewPaginationResponse(paginationReq.Page, normalizedLimit, total)

	return users, pagination, nil
}

// AuthenticateUser authenticates a user.
func (uc *userUsecase) AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: invalid email or password")
	}

	if !user.CheckPassword(password) {
		return nil, fmt.Errorf("authentication failed: invalid email or password")
	}

	user.Password = ""
	return user, nil
}
