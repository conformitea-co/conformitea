package user

import "github.com/google/uuid"

type UserService struct {
	repository UserRepository
}

func Initialize(r UserRepository) *UserService {
	return &UserService{
		repository: r,
	}
}

func (s *UserService) GetUserByID(id uuid.UUID) (*User, error) {
	return s.repository.GetUserByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*User, error) {
	return s.repository.GetUserByEmail(email)
}
