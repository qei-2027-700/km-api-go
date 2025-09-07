package repository

import "km-api-go/internal/domain"

type UserRepository interface {
	GetAll() ([]domain.User, error)
	GetByID(id uint) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id uint) error
	Exists(id uint) (bool, error)
	ExistsByEmail(email string) (bool, error)
	Count() (int64, error)
	GetPaginated(offset, limit int) ([]domain.User, error)
}
