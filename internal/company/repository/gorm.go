package repository

import (
	"errors"
	"fmt"
	"strings"

	"km-api-go/internal/domain"

	"gorm.io/gorm"
)

// companyRepository GORM実装
type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) GetAll() ([]domain.Company, error) {
	var companies []domain.Company

	if err := r.db.Find(&companies).Error; err != nil {
		return nil, fmt.Errorf("failed to get all companies: %w", err)
	}

	return companies, nil
}

func (r *companyRepository) GetByID(id uint) (*domain.Company, error) {
	var c domain.Company

	if err := r.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("company with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get company by id %d: %w", id, err)
	}

	return &c, nil
}

func (r *companyRepository) GetByEmail(email string) (*domain.Company, error) {
	var c domain.Company

	if err := r.db.Where("email = ?", email).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("company with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get company by email %s: %w", email, err)
	}

	return &c, nil
}

func (r *companyRepository) Create(c *domain.Company) error {
	// メール重複チェック
	exists, err := r.ExistsByEmail(c.Email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return fmt.Errorf("company with email %s already exists", c.Email)
	}

	if err := r.db.Create(c).Error; err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}

	return nil
}

func (r *companyRepository) Update(c *domain.Company) error {
	// 会社存在チェック
	exists, err := r.Exists(c.ID)
	if err != nil {
		return fmt.Errorf("failed to check company existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("company with id %d not found", c.ID)
	}

	// メール重複チェック（自分以外）
	var existingCompany domain.Company
	if err := r.db.Where("email = ? AND id != ?", c.Email, c.ID).First(&existingCompany).Error; err == nil {
		return fmt.Errorf("company with email %s already exists", c.Email)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	}

	if err := r.db.Save(c).Error; err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	return nil
}

func (r *companyRepository) Delete(id uint) error {
	// 会社存在チェック
	exists, err := r.Exists(id)
	if err != nil {
		return fmt.Errorf("failed to check company existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("company with id %d not found", id)
	}

	if err := r.db.Delete(&domain.Company{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete company with id %d: %w", id, err)
	}

	return nil
}

func (r *companyRepository) Exists(id uint) (bool, error) {
	var count int64

	if err := r.db.Model(&domain.Company{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check company existence: %w", err)
	}

	return count > 0, nil
}

func (r *companyRepository) ExistsByEmail(email string) (bool, error) {
	var count int64

	if err := r.db.Model(&domain.Company{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check company email existence: %w", err)
	}

	return count > 0, nil
}

func (r *companyRepository) Count() (int64, error) {
	var count int64

	if err := r.db.Model(&domain.Company{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count companies: %w", err)
	}

	return count, nil
}

// GetPaginated ページネーション付きで会社を取得
func (r *companyRepository) GetPaginated(offset, limit int) ([]domain.Company, error) {
	var companies []domain.Company

	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&companies).Error; err != nil {
		return nil, fmt.Errorf("failed to get paginated companies: %w", err)
	}

	return companies, nil
}

// SearchByName 名前で会社を検索
func (r *companyRepository) SearchByName(name string) ([]domain.Company, error) {
	var companies []domain.Company
	searchTerm := "%" + strings.ToLower(name) + "%"

	if err := r.db.Where("LOWER(name) LIKE ?", searchTerm).Find(&companies).Error; err != nil {
		return nil, fmt.Errorf("failed to search companies by name: %w", err)
	}

	return companies, nil
}

// companyUserRepository ユーザー-会社関係リポジトリ GORM実装
type companyUserRepository struct {
	db *gorm.DB
}

// NewCompanyUserRepository ユーザー-会社関係リポジトリのコンストラクタ
func NewCompanyUserRepository(db *gorm.DB) CompanyUserRepository {
	return &companyUserRepository{db: db}
}

func (r *companyUserRepository) GetUsersByCompanyID(companyID uint) ([]domain.CompanyUser, error) {
	var companyUsers []domain.CompanyUser

	if err := r.db.Where("company_id = ?", companyID).Find(&companyUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to get users by company id %d: %w", companyID, err)
	}

	return companyUsers, nil
}

func (r *companyUserRepository) GetCompaniesByUserID(userID uint) ([]domain.CompanyUser, error) {
	var companyUsers []domain.CompanyUser

	if err := r.db.Where("user_id = ?", userID).Find(&companyUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to get companies by user id %d: %w", userID, err)
	}

	return companyUsers, nil
}

func (r *companyUserRepository) Create(companyUser *domain.CompanyUser) error {
	// 既存関係チェック
	exists, err := r.Exists(companyUser.UserID, companyUser.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to check relation existence: %w", err)
	}
	if exists {
		return fmt.Errorf("relation between user %d and company %d already exists", companyUser.UserID, companyUser.CompanyID)
	}

	// 作成実行
	if err := r.db.Create(companyUser).Error; err != nil {
		return fmt.Errorf("failed to create company-user relation: %w", err)
	}

	return nil
}

func (r *companyUserRepository) Update(companyUser *domain.CompanyUser) error {
	// 関係存在チェック
	exists, err := r.Exists(companyUser.UserID, companyUser.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to check relation existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("relation between user %d and company %d not found", companyUser.UserID, companyUser.CompanyID)
	}

	if err := r.db.Where("user_id = ? AND company_id = ?", companyUser.UserID, companyUser.CompanyID).Updates(companyUser).Error; err != nil {
		return fmt.Errorf("failed to update company-user relation: %w", err)
	}

	return nil
}

func (r *companyUserRepository) Delete(userID, companyID uint) error {
	// 関係存在チェック
	exists, err := r.Exists(userID, companyID)
	if err != nil {
		return fmt.Errorf("failed to check relation existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("relation between user %d and company %d not found", userID, companyID)
	}

	if err := r.db.Where("user_id = ? AND company_id = ?", userID, companyID).Delete(&domain.CompanyUser{}).Error; err != nil {
		return fmt.Errorf("failed to delete company-user relation: %w", err)
	}

	return nil
}

// GetRelation ユーザー-会社関係を取得
func (r *companyUserRepository) GetRelation(userID, companyID uint) (*domain.CompanyUser, error) {
	var companyUser domain.CompanyUser

	if err := r.db.Where("user_id = ? AND company_id = ?", userID, companyID).First(&companyUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("relation between user %d and company %d not found", userID, companyID)
		}
		return nil, fmt.Errorf("failed to get relation: %w", err)
	}

	return &companyUser, nil
}

// Exists ユーザー-会社関係の存在確認
func (r *companyUserRepository) Exists(userID, companyID uint) (bool, error) {
	var count int64

	if err := r.db.Model(&domain.CompanyUser{}).Where("user_id = ? AND company_id = ?", userID, companyID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check relation existence: %w", err)
	}

	return count > 0, nil
}
