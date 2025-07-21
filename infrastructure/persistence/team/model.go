package team

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null"`
	Name           string    `gorm:"type:text;not null"`
}

func (u *Team) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID, _ = uuid.NewV7()
	return
}
