package organization

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationRepository interface {
	GetOrganizationByID(DB *gorm.DB, id uuid.UUID) (Organization, error)
}
