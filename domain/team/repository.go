package team

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository interface {
	GetTeamByID(DB *gorm.DB, id uuid.UUID) (Team, error)
}
