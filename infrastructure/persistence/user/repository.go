package user

import (
	domain "conformitea/domain/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct{}

func (u *UserRepository) GetUserByEmail(DB *gorm.DB, email string) (domain.User, error) {
	var user User

	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (u *UserRepository) GetUserByID(DB *gorm.DB, id uuid.UUID) (domain.User, error) {
	var user User

	if err := DB.Where("id = ?", id).First(&user).Error; err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}
