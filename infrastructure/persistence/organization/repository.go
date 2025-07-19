package organization

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindByID(db *gorm.DB, id uuid.UUID) (*Organization, error) {
	var organization Organization

	if err := db.Where("id = ?", id).First(&organization).Error; err != nil {
		return nil, err
	}

	return &organization, nil
}
