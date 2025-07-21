package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	repository UserRepository
}

func Initialize(r UserRepository) *UserService {
	return &UserService{
		repository: r,
	}
}

func (s *UserService) GetUserByID(DB *gorm.DB, id uuid.UUID) (User, error) {
	return s.repository.GetUserByID(DB, id)
}

func (s *UserService) GetUserByEmail(DB *gorm.DB, email string) (User, error) {
	return s.repository.GetUserByEmail(DB, email)
}
