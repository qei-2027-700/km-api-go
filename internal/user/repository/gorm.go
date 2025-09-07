package repository

import (
	"errors"
	"fmt"

	"km-api-go/internal/domain"

	"gorm.io/gorm"
)

// userRepository GORM実装
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository ユーザーリポジトリのコンストラクタ
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll() ([]domain.User, error) {
	var users []domain.User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	return users, nil
}

func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	var u domain.User

	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user by id %d: %w", id, err)
	}

	return &u, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var u domain.User

	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user by email %s: %w", email, err)
	}

	return &u, nil
}

func (r *userRepository) Create(u *domain.User) error {
	// メール重複チェック
	exists, err := r.ExistsByEmail(u.Email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return fmt.Errorf("user with email %s already exists", u.Email)
	}

	// パスワードハッシュ化（ドメインモデルのBeforeCreateフックで実行される）
	if err := r.db.Create(u).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) Update(u *domain.User) error {
	// ユーザー存在チェック
	exists, err := r.Exists(u.ID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user with id %d not found", u.ID)
	}

	// メール重複チェック（自分以外）
	var existingUser domain.User
	if err := r.db.Where("email = ? AND id != ?", u.Email, u.ID).First(&existingUser).Error; err == nil {
		return fmt.Errorf("user with email %s already exists", u.Email)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	}

	// 更新実行
	if err := r.db.Save(u).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *userRepository) Delete(id uint) error {
	// ユーザー存在チェック
	exists, err := r.Exists(id)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user with id %d not found", id)
	}

	// 削除実行
	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user with id %d: %w", id, err)
	}

	return nil
}

func (r *userRepository) Exists(id uint) (bool, error) {
	var count int64

	if err := r.db.Model(&domain.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return count > 0, nil
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64

	if err := r.db.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user email existence: %w", err)
	}

	return count > 0, nil
}

func (r *userRepository) Count() (int64, error) {
	var count int64

	if err := r.db.Model(&domain.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// GetPaginated ページネーション付きでユーザーを取得
func (r *userRepository) GetPaginated(offset, limit int) ([]domain.User, error) {
	var users []domain.User

	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get paginated users: %w", err)
	}

	return users, nil
}
