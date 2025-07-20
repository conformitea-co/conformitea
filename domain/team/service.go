package team

import "github.com/google/uuid"

type TeamService struct {
	repository TeamRepository
}

func Initialize(r TeamRepository) *TeamService {
	return &TeamService{
		repository: r,
	}
}

func (s *TeamService) GetTeamByID(id uuid.UUID) (*Team, error) {
	return s.repository.GetTeamByID(id)
}
