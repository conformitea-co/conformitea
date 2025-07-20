package team

import "github.com/google/uuid"

type TeamRepository interface {
	GetTeamByID(id uuid.UUID) (*Team, error)
}
