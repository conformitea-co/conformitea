package user

import (
	"gorm.io/gorm"
)

type UserRepository struct{}

func (r *UserRepository) GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
