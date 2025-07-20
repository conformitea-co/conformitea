package user

import "github.com/google/uuid"

type UserRepository interface {
	GetUserByID(id uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
}
