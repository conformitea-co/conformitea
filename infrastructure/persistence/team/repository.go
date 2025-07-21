package team

import (
	domain "conformitea/domain/team"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository struct{}

func (t *TeamRepository) GetTeamByID(DB *gorm.DB, id uuid.UUID) (domain.Team, error) {
	var team Team

	if err := DB.Where("id = ?", id).First(&team).Error; err != nil {
		return domain.Team{}, err
	}

	return domain.Team{
		ID:   team.ID,
		Name: team.Name,
	}, nil
}
