package organization

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationService struct {
	repository OrganizationRepository
}

func Initialize(r OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		repository: r,
	}
}

func (s *OrganizationService) GetOrganizationByID(DB *gorm.DB, id uuid.UUID) (Organization, error) {
	return s.repository.GetOrganizationByID(DB, id)
}
