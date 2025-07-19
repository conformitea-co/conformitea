package organization

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"type:text;not null"`
}

func (o *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID, _ = uuid.NewV7()
	return
}
