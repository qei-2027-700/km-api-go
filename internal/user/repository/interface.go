package repository

import "km-api-go/internal/domain"

// UserRepository ユーザーリポジトリインターフェース
type UserRepository interface {
	// GetAll 全ユーザーを取得
	GetAll() ([]domain.User, error)
	
	// GetByID IDでユーザーを取得
	GetByID(id uint) (*domain.User, error)
	
	// GetByEmail メールアドレスでユーザーを取得
	GetByEmail(email string) (*domain.User, error)
	
	// Create ユーザーを作成
	Create(user *domain.User) error
	
	// Update ユーザーを更新
	Update(user *domain.User) error
	
	// Delete IDでユーザーを削除
	Delete(id uint) error
	
	// Exists IDでユーザーの存在確認
	Exists(id uint) (bool, error)
	
	// ExistsByEmail メールアドレスでユーザーの存在確認
	ExistsByEmail(email string) (bool, error)
	
	// Count 総ユーザー数を取得
	Count() (int64, error)
	
	// GetPaginated ページネーション付きでユーザーを取得
	GetPaginated(offset, limit int) ([]domain.User, error)
}