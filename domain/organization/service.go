package organization

import "github.com/google/uuid"

type OrganizationService struct {
	repository OrganizationRepository
}

func Initialize(r OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		repository: r,
	}
}

func (s *OrganizationService) GetOrganizationByID(id uuid.UUID) (*Organization, error) {
	return s.repository.GetOrganizationByID(id)
}
