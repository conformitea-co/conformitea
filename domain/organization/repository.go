package organization

import "github.com/google/uuid"

type OrganizationRepository interface {
	GetOrganizationByID(id uuid.UUID) (*Organization, error)
}
