package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email string    `gorm:"type:text;not null;unique"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID, _ = uuid.NewV7()
	return
}
