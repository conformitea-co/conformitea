package user

import (
	"time"

	"conformitea/infrastructure/persistence/organization"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID                   `gorm:"type:uuid;primaryKey"`
	Email         string                      `gorm:"type:text;not null;unique"`
	FirstName     string                      `gorm:"type:text"`
	LastName      string                      `gorm:"type:text"`
	CreatedAt     time.Time                   `gorm:"autoCreateTime"`
	UpdatedAt     time.Time                   `gorm:"autoUpdateTime"`
	Organizations []organization.Organization `gorm:"many2many:user_organizations;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID, _ = uuid.NewV7()
	return
}
