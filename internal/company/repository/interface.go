package repository

import "km-api-go/internal/domain"

// CompanyRepository 会社リポジトリインターフェース
type CompanyRepository interface {
	// GetAll 全会社を取得
	GetAll() ([]domain.Company, error)
	
	// GetByID IDで会社を取得
	GetByID(id uint) (*domain.Company, error)
	
	// GetByEmail メールアドレスで会社を取得
	GetByEmail(email string) (*domain.Company, error)
	
	// Create 会社を作成
	Create(company *domain.Company) error
	
	// Update 会社を更新
	Update(company *domain.Company) error
	
	// Delete IDで会社を削除
	Delete(id uint) error
	
	// Exists IDで会社の存在確認
	Exists(id uint) (bool, error)
	
	// ExistsByEmail メールアドレスで会社の存在確認
	ExistsByEmail(email string) (bool, error)
	
	// Count 総会社数を取得
	Count() (int64, error)
	
	// GetPaginated ページネーション付きで会社を取得
	GetPaginated(offset, limit int) ([]domain.Company, error)
	
	// SearchByName 名前で会社を検索
	SearchByName(name string) ([]domain.Company, error)
}

// CompanyUserRepository ユーザー-会社関係リポジトリインターフェース
type CompanyUserRepository interface {
	// GetUsersByCompanyID 会社IDでユーザー一覧を取得
	GetUsersByCompanyID(companyID uint) ([]domain.CompanyUser, error)
	
	// GetCompaniesByUserID ユーザーIDで会社一覧を取得
	GetCompaniesByUserID(userID uint) ([]domain.CompanyUser, error)
	
	// Create ユーザー-会社関係を作成
	Create(companyUser *domain.CompanyUser) error
	
	// Update ユーザー-会社関係を更新
	Update(companyUser *domain.CompanyUser) error
	
	// Delete ユーザー-会社関係を削除
	Delete(userID, companyID uint) error
	
	// GetRelation ユーザー-会社関係を取得
	GetRelation(userID, companyID uint) (*domain.CompanyUser, error)
	
	// Exists ユーザー-会社関係の存在確認
	Exists(userID, companyID uint) (bool, error)
}