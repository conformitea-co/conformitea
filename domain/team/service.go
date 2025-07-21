package team

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamService struct {
	repository TeamRepository
}

func Initialize(r TeamRepository) *TeamService {
	return &TeamService{
		repository: r,
	}
}

func (s *TeamService) GetTeamByID(DB *gorm.DB, id uuid.UUID) (Team, error) {
	return s.repository.GetTeamByID(DB, id)
}
