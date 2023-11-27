package repositories

import (
	"github.com/sarrooo/go-clean/internal/models"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(user *models.User) (err error)
	GetByEmail(email string) (user *models.User, err error)
	UpdateColumns(user *gorm.Model) (err error)
}

type UserRepository struct {
	DB *gorm.DB
}

func (rpt *UserRepository) Create(user *models.User) (err error) {
	return rpt.DB.Create(user).Error
}

// GetByEmail returns user by email
// If user not found, returns nil
// If error occurred, returns error
func (rpt *UserRepository) GetByEmail(email string) (user *models.User, err error) {
	user = &models.User{}
	err = rpt.DB.Where("email = ?", email).Limit(1).Find(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (rpt *UserRepository) UpdateColumns(user *gorm.Model) (err error) {
	return rpt.DB.Model(user).Updates(user).Error
}
