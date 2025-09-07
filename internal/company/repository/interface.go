package repository

import "km-api-go/internal/domain"

type CompanyRepository interface {
	GetAll() ([]domain.Company, error)
	GetByID(id uint) (*domain.Company, error)
	GetByEmail(email string) (*domain.Company, error)
	Create(company *domain.Company) error
	Update(company *domain.Company) error
	Delete(id uint) error
	Exists(id uint) (bool, error)
	ExistsByEmail(email string) (bool, error)
	Count() (int64, error)
	GetPaginated(offset, limit int) ([]domain.Company, error)
	SearchByName(name string) ([]domain.Company, error)
}

// ユーザー-会社関係リポジトリインターフェース
type CompanyUserRepository interface {
	GetUsersByCompanyID(companyID uint) ([]domain.CompanyUser, error)
	GetCompaniesByUserID(userID uint) ([]domain.CompanyUser, error)
	Create(companyUser *domain.CompanyUser) error
	Update(companyUser *domain.CompanyUser) error
	Delete(userID, companyID uint) error
	GetRelation(userID, companyID uint) (*domain.CompanyUser, error)
	Exists(userID, companyID uint) (bool, error)
}
