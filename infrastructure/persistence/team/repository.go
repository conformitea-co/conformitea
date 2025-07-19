package team

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindByID(db *gorm.DB, id uuid.UUID) (*Team, error) {
	var team Team

	if err := db.Where("id = ?", id).First(&team).Error; err != nil {
		return nil, err
	}

	return &team, nil
}
