package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(DB *gorm.DB, id uuid.UUID) (User, error)
	GetUserByEmail(DB *gorm.DB, email string) (User, error)
}
